package valueobject

import (
	"errors"
	"time"
)

var (
	ErrDisplayNameTooShort = errors.New("display name must be at least 2 characters")
	ErrDisplayNameTooLong  = errors.New("display name must be at most 50 characters")
	ErrBioTooLong          = errors.New("bio must be at most 500 characters")
	ErrInvalidBirthDate    = errors.New("invalid birth date")
	ErrTooYoung            = errors.New("must be at least 18 years old")
)

type Profile struct {
	DisplayName string
	Bio         string
	BirthDate   time.Time
	Gender      Gender
	Interests   []string
	Photos      []PhotoURL
}

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type PhotoURL string

func NewProfile(displayName, bio string, birthDate time.Time, gender Gender) (Profile, error) {
	if len(displayName) < 2 {
		return Profile{}, ErrDisplayNameTooShort
	}
	if len(displayName) > 50 {
		return Profile{}, ErrDisplayNameTooLong
	}
	if len(bio) > 500 {
		return Profile{}, ErrBioTooLong
	}

	age := calculateAge(birthDate)
	if age < 18 {
		return Profile{}, ErrTooYoung
	}

	return Profile{
		DisplayName: displayName,
		Bio:         bio,
		BirthDate:   birthDate,
		Gender:      gender,
		Interests:   make([]string, 0),
		Photos:      make([]PhotoURL, 0),
	}, nil
}

func calculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	return age
}

func (p Profile) Age() int {
	return calculateAge(p.BirthDate)
}
