package primitive

type Permission int

const (
	Read Permission = iota + 1
	Create
	Update
	Delete
)

type AppID string

const (
	InternalAppID AppID = "internal"
	WebAppID      AppID = "web"
	MobileAppID   AppID = "mobile"
)

func (a AppID) String() string {
	switch a {
	case InternalAppID:
		return "internal"
	case WebAppID:
		return "web"
	case MobileAppID:
		return "mobile"
	default:
		return "UNSPECIFIED"
	}
}
