package application

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/domain/aggregate"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/domain/commands"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/domain/repository"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/domain/valueobject"
)

var (
	ErrMemberAlreadyExists  = errors.New("member with this email already exists")
	ErrInvalidCredentials   = errors.New("invalid email or password")
)

type MemberService struct {
	repo       repository.MemberRepository
	eventStore EventStore
}

type EventStore interface {
	Append(ctx context.Context, aggregateID string, events []any) error
	Load(ctx context.Context, aggregateID string) ([]any, error)
}

func NewMemberService(repo repository.MemberRepository, eventStore EventStore) *MemberService {
	return &MemberService{
		repo:       repo,
		eventStore: eventStore,
	}
}

func (s *MemberService) RegisterMember(ctx context.Context, cmd commands.RegisterMember) (*aggregate.Member, error) {
	existing, _ := s.repo.GetByEmail(ctx, cmd.Email)
	if existing != nil {
		return nil, ErrMemberAlreadyExists
	}

	email, err := valueobject.NewEmail(cmd.Email)
	if err != nil {
		return nil, err
	}

	member, err := aggregate.NewMember(cmd.MemberID, email)
	if err != nil {
		return nil, err
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SaveWithPassword(ctx, member, string(passwordHash)); err != nil {
		return nil, err
	}

	return member, nil
}

// AuthenticateMember verifies the email and password, returning the member if valid.
func (s *MemberService) AuthenticateMember(ctx context.Context, email, password string) (*aggregate.Member, error) {
	member, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if err == aggregate.ErrMemberNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	passwordHash, err := s.repo.GetPasswordHash(ctx, member.ID())
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return member, nil
}

func (s *MemberService) UpdateProfile(ctx context.Context, cmd commands.UpdateProfile) error {
	member, err := s.repo.GetByID(ctx, cmd.MemberID)
	if err != nil {
		return err
	}

	profile, err := valueobject.NewProfile(
		cmd.DisplayName,
		cmd.Bio,
		cmd.BirthDate,
		valueobject.Gender(cmd.Gender),
	)
	if err != nil {
		return err
	}

	if err := member.UpdateProfile(profile); err != nil {
		return err
	}

	return s.repo.Save(ctx, member)
}

func (s *MemberService) ActivateMember(ctx context.Context, cmd commands.ActivateMember) error {
	member, err := s.repo.GetByID(ctx, cmd.MemberID)
	if err != nil {
		return err
	}

	if err := member.Activate(); err != nil {
		return err
	}

	return s.repo.Save(ctx, member)
}

func (s *MemberService) GetMember(ctx context.Context, memberID string) (*aggregate.Member, error) {
	return s.repo.GetByID(ctx, memberID)
}
