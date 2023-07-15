package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type config struct {
	httpPort string
	pgUrl    string

	appName string

	accessTokenExpireTime int64
	accessTokenSecret     string
}

// defaultConfig create the app config with the its default value
func defaultConfig() config {
	return config{
		httpPort: ":3000",
		pgUrl:    "postgresql://postgres:password2@127.0.0.1:5432/hris",

		appName: "HRIS v1.0.0",

		accessTokenExpireTime: 86400,
		accessTokenSecret:     "supersecrettoken",
	}
}

// parseConfig read from the os env and then overwrite the
// existing default config
func parseConfig() config {
	cfg := defaultConfig()
	var err error

	if val, ok := os.LookupEnv("HTTP_PORT"); ok {
		cfg.httpPort = fmt.Sprintf(":%s", val)
	}

	if val, ok := os.LookupEnv("PG_URL"); ok {
		cfg.pgUrl = val
	}

	if val, ok := os.LookupEnv("APP_NAME"); ok {
		cfg.appName = val
	}

	if val, ok := os.LookupEnv("ACCESS_TOKEN_EXPIRE_TIME"); ok {
		cfg.accessTokenExpireTime, err = strconv.ParseInt(val, 10, 64)
		if err != nil {
			log.Fatal("invalid ACCESS_TOKEN_EXPIRE_TIME value")
		}
	} else {
		log.Fatal("ACCESS_TOKEN_EXPIRE_TIME is not set")
	}

	if val, ok := os.LookupEnv("ACCESS_TOKEN_SECRET"); ok {
		cfg.accessTokenSecret = val
	} else {
		log.Fatal("ACCESS_TOKEN_SECRET is not set")
	}

	if val, ok := os.LookupEnv("OAUTH_CLIENT_ID"); ok {
		cfg.accessTokenSecret = val
	} else {
		log.Fatal("OAUTH_CLIENT_ID is not set")
	}

	if val, ok := os.LookupEnv("OAUTH_CLIENT_SECRET"); ok {
		cfg.accessTokenSecret = val
	} else {
		log.Fatal("OAUTH_CLIENT_SECRET is not set")
	}

	if val, ok := os.LookupEnv("OAUTH_REDIRECT_URL"); ok {
		cfg.accessTokenSecret = val
	} else {
		log.Fatal("OAUTH_REDIRECT_URL is not set")
	}

	if val, ok := os.LookupEnv("OAUTH_STATE"); ok {
		cfg.accessTokenSecret = val
	} else {
		log.Fatal("OAUTH_STATE is not set")
	}

	// TODO: validate redis env
	// TODO: validate smtp env
	// TODO: validate web base url
	// TODO: validate internal web base url

	return cfg
}
