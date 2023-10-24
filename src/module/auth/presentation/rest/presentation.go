package rest

import (
	"hris/module/auth/internal"
	"hris/module/auth/mobile"
)

type AuthPresentation struct {
	Internal *internal.Internal
	Mobile   *mobile.Mobile
}
