package service

import (
	"context"
	"errors"
	"hroost/module/shared/primitive"
	"net/http"

	provinceRepo "hroost/module/region/repo/province"

	"github.com/jackc/pgx/v5"
)

type FindAllProvince struct {
	ID   string
	Name string
}

type FindAllProvinceOut struct {
	primitive.CommonResult

	Provinces []FindAllProvince
}

func (s *RegionService) FindAllProvince(ctx context.Context) (out FindAllProvinceOut) {
	// find all province
	provinces, err := s.ProvinceRepo.FindAllProvince(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusOK, "success")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error")
			return
		}
	}

	// map the response
	s.mapFindAllProvince(provinces, &out)

	out.SetResponse(http.StatusOK, "success")
	return
}

func (s *RegionService) mapFindAllProvince(in []provinceRepo.FindAllProvinceOut, out *FindAllProvinceOut) {
	for _, province := range in {
		var p FindAllProvince
		p.ID = province.ID
		p.Name = province.Name

		out.Provinces = append(out.Provinces, p)
	}
}
