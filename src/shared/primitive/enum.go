package primitive

type AppID string

const (
	CentralAppID AppID = "central"
	TenantAppID  AppID = "tenant"
	MobileAppID  AppID = "mobile"
)

func (a AppID) String() string {
	switch a {
	case CentralAppID:
		return "central"
	case TenantAppID:
		return "tenant"
	case MobileAppID:
		return "mobile"
	default:
		return "UNSPECIFIED"
	}
}
