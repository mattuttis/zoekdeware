package application

import (
	"context"
	"errors"

	"github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/domain/aggregate"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/domain/commands"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/domain/repository"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/domain/valueobject"
)

var (
	ErrMemberAlreadyExists = errors.New("member with this email already exists")
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

	if err := s.repo.Save(ctx, member); err != nil {
		return nil, err
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
