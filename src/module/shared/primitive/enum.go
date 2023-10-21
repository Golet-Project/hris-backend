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
	TenantAppID   AppID = "tenant"
	MobileAppID   AppID = "mobile"
)

func (a AppID) String() string {
	switch a {
	case InternalAppID:
		return "internal"
	case TenantAppID:
		return "tenant"
	case MobileAppID:
		return "mobile"
	default:
		return "UNSPECIFIED"
	}
}
