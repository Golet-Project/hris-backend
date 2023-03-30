package entities

import "time"

type User struct {
	ID             int64
	Uid            string
	Email          string
	Password       string
	FullName       string
	BirthDate      string
	ProfilePicture string
	Address        string
	ProvinceId     string
	RegencyId      string
	DistrictId     string
	VillageId      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
