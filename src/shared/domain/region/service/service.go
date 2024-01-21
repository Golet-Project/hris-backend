package service

import (
	"fmt"
	"hroost/shared/domain/region/db/district"
	"hroost/shared/domain/region/db/province"
	"hroost/shared/domain/region/db/regency"
	"hroost/shared/domain/region/db/village"
)

type Config struct {
	ProvinceDb *province.Db
	RegencyDb  *regency.Db
	DistrictDb *district.Db
	VillageDb  *village.Db
}

type Service struct {
	provinceDb *province.Db
	regencyDb  *regency.Db
	districtDb *district.Db
	villageDb  *village.Db
}

func New(cfg *Config) (*Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config required")
	}
	if cfg.ProvinceDb == nil {
		return nil, fmt.Errorf("provinceDb required")
	}
	if cfg.RegencyDb == nil {
		return nil, fmt.Errorf("regencyDb required")
	}
	if cfg.DistrictDb == nil {
		return nil, fmt.Errorf("districtDb required")
	}
	if cfg.VillageDb == nil {
		return nil, fmt.Errorf("villageDb required")
	}

	return &Service{
		provinceDb: cfg.ProvinceDb,
		regencyDb:  cfg.RegencyDb,
		districtDb: cfg.DistrictDb,
		villageDb:  cfg.VillageDb,
	}, nil
}
