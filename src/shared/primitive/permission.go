package primitive

type Permission int

const (
	Read Permission = iota + 1
	Create
	Update
	Delete
)
