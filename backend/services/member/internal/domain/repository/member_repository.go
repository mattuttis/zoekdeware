package repository

import (
	"context"

	"github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/domain/aggregate"
)

type MemberRepository interface {
	Save(ctx context.Context, member *aggregate.Member) error
	GetByID(ctx context.Context, id string) (*aggregate.Member, error)
	GetByEmail(ctx context.Context, email string) (*aggregate.Member, error)
}
