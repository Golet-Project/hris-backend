package primitive

import (
	"errors"
	"time"
)

var ErrInvalidTimezone = errors.New("invalid timezone")

type Timezone int

const (
	_ Timezone = iota + 6
	WIB
	WITA
	WIT
)

func (t Timezone) Value() int {
	return int(t)
}

func (t Timezone) Valid() bool {
	switch t {
	case WIB, WITA, WIT:
		return true
	default:
		return false
	}
}

func (t Timezone) Now() (time.Time, error) {
	var loc *time.Location
	switch t {
	case WIB:
		loc = time.FixedZone("WIB", 7*60*60)
	case WITA:
		loc = time.FixedZone("WITA", 8*60*60)
	case WIT:
		loc = time.FixedZone("WIT", 9*60*60)
	default:
		return time.Time{}, ErrInvalidTimezone
	}

	return time.Now().In(loc), nil
}