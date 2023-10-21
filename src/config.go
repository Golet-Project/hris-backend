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

	redisHost     string
	redisPort     string
	redisPassword string
	redisDB       string
	redisTaskDB   string
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

	// pg database
	pgHost, ok := os.LookupEnv("PG_HOST")
	if !ok {
		log.Fatal("PG_HOST is not set")
	} else {
		if len(pgHost) == 0 {
			log.Fatal("PG_HOST can not be empty")
		}
	}

	pgPort, ok := os.LookupEnv("PG_PORT")
	if !ok {
		log.Fatal("PG_PORT is not set")
	} else {
		if len(pgPort) == 0 {
			log.Fatal("PG_PORT can not be empty")
		}
	}

	pgUser, ok := os.LookupEnv("PG_USER")
	if !ok {
		log.Fatal("PG_USER is not set")
	} else {
		if len(pgUser) == 0 {
			log.Fatal("PG_USER can not be empty")
		}
	}
	pgPassword, ok := os.LookupEnv("PG_PASSWORD")
	if !ok {
		log.Fatal("PG_PASSWORD is not set")
	} else {
		if len(pgPassword) == 0 {
			log.Fatal("PG_PASSWORD can not be empty")
		}
	}

	pgMasterDB, ok := os.LookupEnv("PG_MASTER_DB")
	if !ok {
		log.Fatal("PG_MASTER_DB is not set")
	} else {
		if len(pgMasterDB) == 0 {
			log.Fatal("PG_MASTER_DB can not be empty")
		}
	}
	cfg.pgUrl = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		pgUser,
		pgPassword,
		pgHost,
		pgPort,
		pgMasterDB,
	)

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

	// REDIS
	if val, ok := os.LookupEnv("REDIS_HOST"); ok {
		cfg.redisHost = val
	} else {
		log.Fatal("REDIS_HOST is not set")
	}

	if val, ok := os.LookupEnv("REDIS_PORT"); ok {
		cfg.redisPort = val
	} else {
		log.Fatal("REDIS_PORT is not set")
	}

	if val, ok := os.LookupEnv("REDIS_PASSWORD"); ok {
		cfg.redisPassword = val
	} else {
		log.Fatal("REDIS_PASSWORD is not set")
	}

	if val, ok := os.LookupEnv("REDIS_TASK_DB"); ok {
		cfg.redisTaskDB = val
	} else {
		log.Fatal("REDIS_TASK_DB is not set")
	}

	// SMTPc
	if val, ok := os.LookupEnv("SMTP_HOST"); ok {
		if len(val) == 0 {
			log.Fatal("SMPTP_HOST can not be empty")
		}
	} else {
		log.Fatal("SMTP_HOST is not set")
	}

	if val, ok := os.LookupEnv("SMTP_PORT"); ok {
		if len(val) == 0 {
			log.Fatal("SMTP_PORT can not be empty")
		}
	} else {
		log.Fatal("SMTP_PORT is not set")
	}

	if val, ok := os.LookupEnv("SMTP_AUTH_USER"); ok {
		if len(val) == 0 {
			log.Fatal("SMTP_AUTH_USER can not be empty")
		}
	} else {
		log.Fatal("SMTP_AUTH_USER is not set")
	}

	if val, ok := os.LookupEnv("SMTP_AUTH_PASSWORD"); ok {
		if len(val) == 0 {
			log.Fatal("SMTP_AUTH_PASSWORD can not set")
		}
	} else {
		log.Fatal("SMTP_AUTH_PASSWORD is not set")
	}

	if val, ok := os.LookupEnv("SMTP_SENDER_NAME"); ok {
		if len(val) == 0 {
			log.Fatal("SMTP_SENDER_NAME can not be empty")
		}
	} else {
		log.Fatal("SMTP_SENDER_NAME is not set")
	}

	// WEB BASE URL
	if val, ok := os.LookupEnv("WEB_BASE_URL"); ok {
		if len(val) == 0 {
			log.Fatal("WEB_BASE_URL can not be empty")
		}
	} else {
		log.Fatal("WEB_BASE_URL is not set")
	}

	// INTERAL WEB BASE URL
	if val, ok := os.LookupEnv("INTERNAL_WEB_BASE_URL"); ok {
		if len(val) == 0 {
			log.Fatal("INTERNAL_WEB_BASE_URL can not be empty")
		}
	} else {
		log.Fatal("INTERAL_WEB_BASE_URL is not set")
	}

	return cfg
}
