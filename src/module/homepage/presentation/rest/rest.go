package rest

import (
	"hroost/module/homepage/mobile"
	"log"
)

type Rest struct {
	mobile *mobile.Mobile
}

func New(mobile *mobile.Mobile) *Rest {
	if mobile == nil {
		log.Fatal("[x] Mobile service required on homepage/presentation/rest package")
	}

	return &Rest{
		mobile: mobile,
	}
}
