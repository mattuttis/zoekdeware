package commands

import "time"

type Command interface {
	CommandType() string
}

type RegisterMember struct {
	MemberID string
	Email    string
	Password string
}

func (c RegisterMember) CommandType() string { return "member.register" }

type UpdateProfile struct {
	MemberID    string
	DisplayName string
	Bio         string
	BirthDate   time.Time
	Gender      string
}

func (c UpdateProfile) CommandType() string { return "member.update_profile" }

type ActivateMember struct {
	MemberID string
}

func (c ActivateMember) CommandType() string { return "member.activate" }

type SuspendMember struct {
	MemberID string
	Reason   string
}

func (c SuspendMember) CommandType() string { return "member.suspend" }
