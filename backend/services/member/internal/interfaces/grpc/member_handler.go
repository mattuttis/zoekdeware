package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	memberv1 "github.com/mattuttis/inetcontrol/zoekdeware/api/proto/member/v1"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/application"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/domain/aggregate"
	"github.com/mattuttis/inetcontrol/zoekdeware/backend/services/member/internal/domain/commands"
)

// MemberHandler implements the gRPC MemberServiceServer interface.
type MemberHandler struct {
	memberv1.UnimplementedMemberServiceServer
	service *application.MemberService
}

// NewMemberHandler creates a new gRPC handler for the member service.
func NewMemberHandler(service *application.MemberService) *MemberHandler {
	return &MemberHandler{service: service}
}

// RegisterMember handles member registration requests.
func (h *MemberHandler) RegisterMember(ctx context.Context, req *memberv1.RegisterMemberRequest) (*memberv1.RegisterMemberResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	cmd := commands.RegisterMember{
		MemberID: uuid.New().String(),
		Email:    req.Email,
		Password: req.Password,
	}

	member, err := h.service.RegisterMember(ctx, cmd)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &memberv1.RegisterMemberResponse{
		Member: toProtoMember(member),
	}, nil
}

// AuthenticateMember handles member authentication requests.
func (h *MemberHandler) AuthenticateMember(ctx context.Context, req *memberv1.AuthenticateMemberRequest) (*memberv1.AuthenticateMemberResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	member, err := h.service.AuthenticateMember(ctx, req.Email, req.Password)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &memberv1.AuthenticateMemberResponse{
		Member: toProtoMember(member),
	}, nil
}

// GetMember retrieves a member by ID.
func (h *MemberHandler) GetMember(ctx context.Context, req *memberv1.GetMemberRequest) (*memberv1.GetMemberResponse, error) {
	if req.MemberId == "" {
		return nil, status.Error(codes.InvalidArgument, "member_id is required")
	}

	member, err := h.service.GetMember(ctx, req.MemberId)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &memberv1.GetMemberResponse{
		Member: toProtoMember(member),
	}, nil
}

// UpdateProfile updates a member's profile.
func (h *MemberHandler) UpdateProfile(ctx context.Context, req *memberv1.UpdateProfileRequest) (*memberv1.UpdateProfileResponse, error) {
	if req.MemberId == "" {
		return nil, status.Error(codes.InvalidArgument, "member_id is required")
	}
	if req.Profile == nil {
		return nil, status.Error(codes.InvalidArgument, "profile is required")
	}

	var birthDate = req.Profile.BirthDate.AsTime()

	cmd := commands.UpdateProfile{
		MemberID:    req.MemberId,
		DisplayName: req.Profile.DisplayName,
		Bio:         req.Profile.Bio,
		BirthDate:   birthDate,
		Gender:      protoGenderToString(req.Profile.Gender),
	}

	if err := h.service.UpdateProfile(ctx, cmd); err != nil {
		return nil, toGRPCError(err)
	}

	member, err := h.service.GetMember(ctx, req.MemberId)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &memberv1.UpdateProfileResponse{
		Member: toProtoMember(member),
	}, nil
}

// ActivateMember activates a pending member.
func (h *MemberHandler) ActivateMember(ctx context.Context, req *memberv1.ActivateMemberRequest) (*memberv1.ActivateMemberResponse, error) {
	if req.MemberId == "" {
		return nil, status.Error(codes.InvalidArgument, "member_id is required")
	}

	cmd := commands.ActivateMember{
		MemberID: req.MemberId,
	}

	if err := h.service.ActivateMember(ctx, cmd); err != nil {
		return nil, toGRPCError(err)
	}

	member, err := h.service.GetMember(ctx, req.MemberId)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &memberv1.ActivateMemberResponse{
		Member: toProtoMember(member),
	}, nil
}

// toProtoMember converts a domain member to a protobuf member.
func toProtoMember(m *aggregate.Member) *memberv1.Member {
	profile := m.Profile()

	photoURLs := make([]string, len(profile.Photos))
	for i, p := range profile.Photos {
		photoURLs[i] = string(p)
	}

	return &memberv1.Member{
		Id:     m.ID(),
		Email:  m.Email().String(),
		Status: toProtoStatus(m.Status()),
		Profile: &memberv1.Profile{
			DisplayName: profile.DisplayName,
			Bio:         profile.Bio,
			BirthDate:   timestamppb.New(profile.BirthDate),
			Gender:      toProtoGender(string(profile.Gender)),
			Interests:   profile.Interests,
			PhotoUrls:   photoURLs,
		},
	}
}

// toProtoStatus converts domain status to protobuf status.
func toProtoStatus(s aggregate.MemberStatus) memberv1.MemberStatus {
	switch s {
	case aggregate.MemberStatusPending:
		return memberv1.MemberStatus_MEMBER_STATUS_PENDING
	case aggregate.MemberStatusActive:
		return memberv1.MemberStatus_MEMBER_STATUS_ACTIVE
	case aggregate.MemberStatusSuspended:
		return memberv1.MemberStatus_MEMBER_STATUS_SUSPENDED
	default:
		return memberv1.MemberStatus_MEMBER_STATUS_UNSPECIFIED
	}
}

// toProtoGender converts domain gender to protobuf gender.
func toProtoGender(g string) memberv1.Gender {
	switch g {
	case "male":
		return memberv1.Gender_GENDER_MALE
	case "female":
		return memberv1.Gender_GENDER_FEMALE
	case "other":
		return memberv1.Gender_GENDER_OTHER
	default:
		return memberv1.Gender_GENDER_UNSPECIFIED
	}
}

// protoGenderToString converts protobuf gender to string.
func protoGenderToString(g memberv1.Gender) string {
	switch g {
	case memberv1.Gender_GENDER_MALE:
		return "male"
	case memberv1.Gender_GENDER_FEMALE:
		return "female"
	case memberv1.Gender_GENDER_OTHER:
		return "other"
	default:
		return ""
	}
}

// toGRPCError converts domain errors to gRPC status errors.
func toGRPCError(err error) error {
	switch err {
	case aggregate.ErrMemberNotFound:
		return status.Error(codes.NotFound, err.Error())
	case aggregate.ErrInvalidEmail:
		return status.Error(codes.InvalidArgument, err.Error())
	case application.ErrMemberAlreadyExists:
		return status.Error(codes.AlreadyExists, err.Error())
	case application.ErrInvalidCredentials:
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
