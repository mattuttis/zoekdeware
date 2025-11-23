package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/domain/aggregate"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/domain/events"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/domain/repository"
)

const aggregateType = "Member"

// PostgresMemberRepository implements repository.MemberRepository using PostgreSQL
// with event sourcing pattern.
type PostgresMemberRepository struct {
	db *sql.DB
}

// NewPostgresMemberRepository creates a new PostgreSQL-backed member repository.
func NewPostgresMemberRepository(db *sql.DB) repository.MemberRepository {
	return &PostgresMemberRepository{db: db}
}

// Save persists all uncommitted events from the member aggregate to the event store
// and updates the read model within a single transaction.
func (r *PostgresMemberRepository) Save(ctx context.Context, member *aggregate.Member) error {
	changes := member.Changes()
	if len(changes) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Get current version for optimistic locking
	currentVersion := member.Version() - len(changes)

	// Save each event to the event store
	for i, event := range changes {
		eventData, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("marshal event: %w", err)
		}

		version := currentVersion + i + 1
		_, err = tx.ExecContext(ctx, `
			INSERT INTO events (aggregate_id, aggregate_type, event_type, event_data, version, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, event.AggregateID(), aggregateType, event.EventType(), eventData, version, event.OccurredAt())

		if err != nil {
			return fmt.Errorf("insert event: %w", err)
		}
	}

	// Update read model
	if err := r.updateReadModel(ctx, tx, member); err != nil {
		return fmt.Errorf("update read model: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	member.ClearChanges()
	return nil
}

// GetByID retrieves a member by rehydrating from the event stream.
func (r *PostgresMemberRepository) GetByID(ctx context.Context, id string) (*aggregate.Member, error) {
	eventStream, err := r.loadEvents(ctx, "aggregate_id = $1", id)
	if err != nil {
		return nil, err
	}

	if len(eventStream) == 0 {
		return nil, aggregate.ErrMemberNotFound
	}

	return aggregate.RehydrateMember(eventStream), nil
}

// GetByEmail retrieves a member by email using the read model for lookup,
// then rehydrates from the event stream.
func (r *PostgresMemberRepository) GetByEmail(ctx context.Context, email string) (*aggregate.Member, error) {
	// Use read model for email lookup
	var id string
	err := r.db.QueryRowContext(ctx, `
		SELECT id FROM members WHERE email = $1
	`, email).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, aggregate.ErrMemberNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query member by email: %w", err)
	}

	return r.GetByID(ctx, id)
}

// loadEvents retrieves events from the event store and deserializes them.
func (r *PostgresMemberRepository) loadEvents(ctx context.Context, where string, args ...any) ([]events.Event, error) {
	query := fmt.Sprintf(`
		SELECT event_type, event_data, version
		FROM events
		WHERE %s AND aggregate_type = $%d
		ORDER BY version ASC
	`, where, len(args)+1)

	args = append(args, aggregateType)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query events: %w", err)
	}
	defer rows.Close()

	var eventStream []events.Event
	for rows.Next() {
		var (
			eventType string
			eventData []byte
			version   int
		)
		if err := rows.Scan(&eventType, &eventData, &version); err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}

		event, err := deserializeEvent(eventType, eventData)
		if err != nil {
			return nil, fmt.Errorf("deserialize event: %w", err)
		}
		eventStream = append(eventStream, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate events: %w", err)
	}

	return eventStream, nil
}

// updateReadModel upserts the member read model from the current aggregate state.
func (r *PostgresMemberRepository) updateReadModel(ctx context.Context, tx *sql.Tx, member *aggregate.Member) error {
	profile := member.Profile()

	// Convert photo URLs to string slice for PostgreSQL array
	photos := make([]string, len(profile.Photos))
	for i, p := range profile.Photos {
		photos[i] = string(p)
	}

	_, err := tx.ExecContext(ctx, `
		INSERT INTO members (id, email, display_name, bio, birth_date, gender, interests, photos, status, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			display_name = EXCLUDED.display_name,
			bio = EXCLUDED.bio,
			birth_date = EXCLUDED.birth_date,
			gender = EXCLUDED.gender,
			interests = EXCLUDED.interests,
			photos = EXCLUDED.photos,
			status = EXCLUDED.status,
			version = EXCLUDED.version,
			updated_at = NOW()
	`,
		member.ID(),
		member.Email().String(),
		nullString(profile.DisplayName),
		nullString(profile.Bio),
		nullTime(profile.BirthDate),
		nullString(string(profile.Gender)),
		pq.Array(profile.Interests),
		pq.Array(photos),
		string(member.Status()),
		member.Version(),
	)

	return err
}

// nullString returns sql.NullString for optional string fields.
func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

// nullTime returns sql.NullTime for optional time fields.
func nullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: t, Valid: true}
}

// deserializeEvent converts raw event data back into typed event structs.
func deserializeEvent(eventType string, data []byte) (events.Event, error) {
	switch eventType {
	case "member.registered":
		var e events.MemberRegistered
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return e, nil

	case "member.profile_updated":
		var e events.ProfileUpdated
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return e, nil

	case "member.activated":
		var e events.MemberActivated
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return e, nil

	case "member.suspended":
		var e events.MemberSuspended
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return e, nil

	default:
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}
}
