package primitive

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
