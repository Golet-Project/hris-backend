package service

import (
	"os"
)

var oauthState = os.Getenv("OAUTH_STATE")
