package service

import (
	"context"
	"errors"
	"hroost/shared/primitive"
	"net/http"

	provinceRepo "hroost/shared/domain/region/db/province"

	"github.com/jackc/pgx/v5"
)

type FindAllProvince struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type FindAllProvinceOut struct {
	primitive.CommonResult

	Provinces []FindAllProvince
}

func (s *Service) FindAllProvince(ctx context.Context) (out FindAllProvinceOut) {
	// find all province
	provinces, err := s.provinceDb.FindAllProvince(ctx)
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

func (s *Service) mapFindAllProvince(in []provinceRepo.FindAllProvinceOut, out *FindAllProvinceOut) {
	for _, province := range in {
		var p FindAllProvince
		p.ID = province.ID
		p.Name = province.Name

		out.Provinces = append(out.Provinces, p)
	}
}
