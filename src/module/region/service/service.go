package service

import "hris/module/region/repo/province"

type RegionService struct {
	ProvinceRepo *province.Repository
}

func NewRegionService(provinceRepo *province.Repository) *RegionService {
	return &RegionService{
		ProvinceRepo: provinceRepo,
	}
}
