package service

import (
	"context"
	"errors"
	"hroost/shared/entities"
	"hroost/shared/primitive"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type FindAllRegencyByProvinceIdIn struct {
	ProvinceId string `query:"province_id"`
}

type FindAllRegencyByProvinceId struct {
	entities.Regency
}

type FindAllRegencyByProvinceIdOut struct {
	primitive.CommonResult

	Regency []FindAllRegencyByProvinceId
}

func (s Service) FindAllRegencyByProvinceId(ctx context.Context, req FindAllRegencyByProvinceIdIn) (out FindAllRegencyByProvinceIdOut) {
	if req.ProvinceId == "" {
		out.SetResponse(http.StatusOK, "success")
		return
	}

	// find all regency by province id
	regencies, err := s.regencyDb.FindAllByProvinceId(ctx, req.ProvinceId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusOK, "success")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error")
			return
		}
	}

	for _, regency := range regencies {
		var r FindAllRegencyByProvinceId
		r.Id = regency.Id
		r.ProvinceId = regency.ProvinceId
		r.Name = regency.Name

		out.Regency = append(out.Regency, r)
	}

	out.SetResponse(http.StatusOK, "success")

	return
}
