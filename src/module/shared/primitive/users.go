package primitive

type UserType string

// Enum for authenticated user type
const (
	UserTypeEmployee UserType = "employee"
	UserTypeRoot     UserType = "root"
)
