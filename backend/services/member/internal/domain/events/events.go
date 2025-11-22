package events

import "time"

type Event interface {
	EventType() string
	AggregateID() string
	OccurredAt() time.Time
}

type MemberRegistered struct {
	MemberID  string
	Email     string
	Timestamp time.Time
}

func (e MemberRegistered) EventType() string    { return "member.registered" }
func (e MemberRegistered) AggregateID() string  { return e.MemberID }
func (e MemberRegistered) OccurredAt() time.Time { return e.Timestamp }

type ProfileUpdated struct {
	MemberID    string
	DisplayName string
	Bio         string
	BirthDate   time.Time
	Timestamp   time.Time
}

func (e ProfileUpdated) EventType() string    { return "member.profile_updated" }
func (e ProfileUpdated) AggregateID() string  { return e.MemberID }
func (e ProfileUpdated) OccurredAt() time.Time { return e.Timestamp }

type MemberActivated struct {
	MemberID  string
	Timestamp time.Time
}

func (e MemberActivated) EventType() string    { return "member.activated" }
func (e MemberActivated) AggregateID() string  { return e.MemberID }
func (e MemberActivated) OccurredAt() time.Time { return e.Timestamp }

type MemberSuspended struct {
	MemberID  string
	Reason    string
	Timestamp time.Time
}

func (e MemberSuspended) EventType() string    { return "member.suspended" }
func (e MemberSuspended) AggregateID() string  { return e.MemberID }
func (e MemberSuspended) OccurredAt() time.Time { return e.Timestamp }
