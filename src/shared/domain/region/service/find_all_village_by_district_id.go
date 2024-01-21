package service

import (
	"context"
	"errors"
	"hroost/shared/entities"
	"hroost/shared/primitive"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type FindAllVillageByDistrictIdIn struct {
	DistrictId string `query:"district_id"`
}

type FindAllVillageByDistrictId struct {
	entities.Village
}

type FindAllVillageByDistrictIdOut struct {
	primitive.CommonResult

	Village []FindAllVillageByDistrictId
}

func (s Service) FindAllVillageByDistrictId(ctx context.Context, req FindAllVillageByDistrictIdIn) (out FindAllVillageByDistrictIdOut) {
	if req.DistrictId == "" {
		out.SetResponse(http.StatusOK, "success")
		return
	}

	// find all villlages by district_id
	villages, err := s.villageDb.FindAllByDistrictId(ctx, req.DistrictId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusOK, "success")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error")
			return
		}
	}

	for _, village := range villages {
		var v FindAllVillageByDistrictId
		v.Id = village.Id
		v.DistrictId = village.DistrictId
		v.Name = village.Name

		out.Village = append(out.Village, v)
	}

	out.SetResponse(http.StatusOK, "success")

	return
}
