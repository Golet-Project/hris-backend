package region

import (
	"hroost/module/region/presentation/rest"
	"hroost/module/region/repo/province"
	"hroost/module/region/service"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Region struct {
	RegionPresenter *rest.RegionPresenter
}

type Dependency struct {
	DB *pgxpool.Pool
}

func InitRegion(d *Dependency) *Region {
	if d.DB == nil {
		log.Fatal("[x] Region package required a database connection")
	}

	// init repo
	provinceRepo := province.Repository{
		DB: d.DB,
	}

	// init service
	regionService := service.NewRegionService(&provinceRepo)

	// init presenter
	regionPresenter := rest.RegionPresenter{
		RegionService: regionService,
	}

	return &Region{
		RegionPresenter: &regionPresenter,
	}
}
