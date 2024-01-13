package service

import "hroost/module/region/repo/province"

type RegionService struct {
	ProvinceRepo *province.Repository
}

func NewRegionService(provinceRepo *province.Repository) *RegionService {
	return &RegionService{
		ProvinceRepo: provinceRepo,
	}
}
