package service

import (
	"context"
	"errors"
	"hroost/shared/entities"
	"hroost/shared/primitive"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type FindAllDistrictByRegencyIdIn struct {
	RegencyId string `query:"regency_id"`
}

type FindAllDistrictByRegencyId struct {
	entities.District
}

type FindAllDistrictByRegencyIdOut struct {
	primitive.CommonResult

	District []FindAllDistrictByRegencyId
}

func (s Service) FindAllDistrictByRegencyId(ctx context.Context, req FindAllDistrictByRegencyIdIn) (out FindAllDistrictByRegencyIdOut) {
	if req.RegencyId == "" {
		out.SetResponse(http.StatusOK, "success")
		return
	}

	// find all district by regency_id
	districts, err := s.districtDb.FindAllByRegencyId(ctx, req.RegencyId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusOK, "success")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error")
			return
		}
	}

	for _, district := range districts {
		var d FindAllDistrictByRegencyId
		d.Id = district.Id
		d.RegencyId = district.RegencyId
		d.Name = district.Name

		out.District = append(out.District, d)
	}

	out.SetResponse(http.StatusOK, "success")

	return
}
