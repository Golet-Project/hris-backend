package rest

import (
	"hris/module/attendance/mobile"
	"log"
)

type AttandancePresentation struct {
	mobile *mobile.Mobile
}

func New(mobile *mobile.Mobile) *AttandancePresentation {
	if mobile == nil {
		log.Fatal("[x] Mobile service required on attendance presentation")
	}

	return &AttandancePresentation{
		mobile: mobile,
	}
}
