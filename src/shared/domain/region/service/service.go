package service

import (
	"fmt"
	"hroost/shared/domain/region/db/province"
)

type Config struct {
	ProvinceDb *province.Db
}

type Service struct {
	provinceDb *province.Db
}

func New(cfg *Config) (*Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config required")
	}
	if cfg.ProvinceDb == nil {
		return nil, fmt.Errorf("provinceDb required")
	}

	return &Service{
		provinceDb: cfg.ProvinceDb,
	}, nil
}
