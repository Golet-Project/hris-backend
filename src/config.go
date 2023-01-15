package main

import "os"

type config struct {
	httpPort string
	pgUrl    string

	appName string

	accessTokenSecret string
}

// defaultConfig create the app config with the its default value
func defaultConfig() config {
	return config{
		httpPort: ":3000",
		pgUrl:    "postgresql://postgres:password2@127.0.0.1:5432/hris",

		appName: "HRIS v1.0.0",

		accessTokenSecret: "supersecrettoken",
	}
}

// parseConfig read from the os env and then overwrite the
// existing default config
func parseConfig() config {
	cfg := defaultConfig()

	if val, ok := os.LookupEnv("HTTP_PORT"); ok {
		cfg.httpPort = val
	}

	if val, ok := os.LookupEnv("PG_URL"); ok {
		cfg.pgUrl = val
	}

	if val, ok := os.LookupEnv("APP_NAME"); ok {
		cfg.appName = val
	}

	if val, ok := os.LookupEnv("ACCESS_TOKEN_SECRET"); ok {
		cfg.accessTokenSecret = val
	}

	return cfg
}
