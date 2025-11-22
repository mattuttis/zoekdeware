package aggregate

import (
	"errors"
	"time"

	"github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/domain/events"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/domain/valueobject"
)

var (
	ErrMemberNotFound     = errors.New("member not found")
	ErrInvalidEmail       = errors.New("invalid email address")
	ErrProfileIncomplete  = errors.New("profile is incomplete")
)

type Member struct {
	id        string
	email     valueobject.Email
	profile   valueobject.Profile
	status    MemberStatus
	createdAt time.Time
	updatedAt time.Time
	version   int

	changes []events.Event
}

type MemberStatus string

const (
	MemberStatusPending  MemberStatus = "pending"
	MemberStatusActive   MemberStatus = "active"
	MemberStatusSuspended MemberStatus = "suspended"
)

func NewMember(id string, email valueobject.Email) (*Member, error) {
	m := &Member{
		id:        id,
		email:     email,
		status:    MemberStatusPending,
		createdAt: time.Now(),
		updatedAt: time.Now(),
		version:   0,
		changes:   make([]events.Event, 0),
	}

	m.raise(events.MemberRegistered{
		MemberID:  id,
		Email:     email.String(),
		Timestamp: m.createdAt,
	})

	return m, nil
}

func (m *Member) ID() string {
	return m.id
}

func (m *Member) Email() valueobject.Email {
	return m.email
}

func (m *Member) Status() MemberStatus {
	return m.status
}

func (m *Member) Version() int {
	return m.version
}

func (m *Member) UpdateProfile(profile valueobject.Profile) error {
	m.profile = profile
	m.updatedAt = time.Now()

	m.raise(events.ProfileUpdated{
		MemberID:    m.id,
		DisplayName: profile.DisplayName,
		Bio:         profile.Bio,
		BirthDate:   profile.BirthDate,
		Timestamp:   m.updatedAt,
	})

	return nil
}

func (m *Member) Activate() error {
	if m.status == MemberStatusActive {
		return nil
	}

	m.status = MemberStatusActive
	m.updatedAt = time.Now()

	m.raise(events.MemberActivated{
		MemberID:  m.id,
		Timestamp: m.updatedAt,
	})

	return nil
}

func (m *Member) raise(event events.Event) {
	m.changes = append(m.changes, event)
}

func (m *Member) Changes() []events.Event {
	return m.changes
}

func (m *Member) ClearChanges() {
	m.changes = make([]events.Event, 0)
}

func (m *Member) Apply(event events.Event) {
	switch e := event.(type) {
	case events.MemberRegistered:
		m.id = e.MemberID
		m.email = valueobject.Email(e.Email)
		m.status = MemberStatusPending
		m.createdAt = e.Timestamp
	case events.ProfileUpdated:
		m.profile = valueobject.Profile{
			DisplayName: e.DisplayName,
			Bio:         e.Bio,
			BirthDate:   e.BirthDate,
		}
		m.updatedAt = e.Timestamp
	case events.MemberActivated:
		m.status = MemberStatusActive
		m.updatedAt = e.Timestamp
	}
	m.version++
}

func RehydrateMember(eventStream []events.Event) *Member {
	m := &Member{
		changes: make([]events.Event, 0),
	}
	for _, event := range eventStream {
		m.Apply(event)
	}
	return m
}
