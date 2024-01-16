package service

import "hroost/domain/shared/region/db/province"

type RegionService struct {
	ProvinceRepo *province.Repository
}

func NewRegionService(provinceRepo *province.Repository) *RegionService {
	return &RegionService{
		ProvinceRepo: provinceRepo,
	}
}
