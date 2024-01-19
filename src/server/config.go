package server

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type config struct {
	appName string

	httpPort string

	pgMasterUser     string
	pgMasterPassword string
	pgMasterHost     string
	pgMasterPort     string
	pgMasterDatabase string

	accessTokenExpireTime int64
	accessTokenSecret     string

	googleOAuthClientId     string
	googleOAuthClientSecret string
	googleOAuthRedirectUrl  string
	googleOAuthState        string

	redisMasterHost     string
	redisMasterPort     string
	redisMasterPassword string
	redisMasterDb       int

	asynqRedisMasterHost     string
	asynqRedisMasterPort     string
	asynqRedisMasterPassword string
	asynqRedisMasterDb       int

	smtpHost         string
	smtpPort         string
	smtpAuthUser     string
	smtpAuthPassword string
	smtpSenderName   string

	webBaseUrl string

	centralWebBaseUrl string
}

func parseConfig() config {
	cfg := config{}

	if val, ok := os.LookupEnv("APP_NAME"); ok && len(val) > 0 {
		cfg.appName = val
	} else {
		log.Fatal("APP_NAME is not set")
	}

	if val, ok := os.LookupEnv("HTTP_PORT"); ok && len(val) > 0 {
		cfg.httpPort = fmt.Sprintf(":%s", val)
	} else {
		log.Fatal("HTTP_PORT")
	}

	if val, ok := os.LookupEnv("PG_MASTER_USER"); ok && len(val) > 0 {
		cfg.pgMasterUser = val
	} else {
		log.Fatal("PG_MASTER_USER is not set")
	}

	if val, ok := os.LookupEnv("PG_MASTER_PASSWORD"); ok && len(val) > 0 {
		cfg.pgMasterPassword = val
	} else {
		log.Fatal("PG_MASTER_PASSWORD is not set")
	}

	if val, ok := os.LookupEnv("PG_MASTER_HOST"); ok && len(val) > 0 {
		cfg.pgMasterHost = val
	} else {
		log.Fatal("PG_MASTER_HOST is not set")
	}

	if val, ok := os.LookupEnv("PG_MASTER_PORT"); ok && len(val) > 0 {
		cfg.pgMasterPort = val
	} else {
		log.Fatal("PG_MASTER_PORT is not set")
	}

	if val, ok := os.LookupEnv("PG_MASTER_DATABASE"); ok && len(val) > 0 {
		cfg.pgMasterDatabase = val
	} else {
		log.Fatal("PG_MASTER_DATABASE is not set")
	}

	if val, ok := os.LookupEnv("ACCESS_TOKEN_EXPIRE_TIME"); ok && len(val) > 0 {
		expTime, err := strconv.ParseInt(val, 10, 64)
		cfg.accessTokenExpireTime = expTime
		if err != nil {
			log.Fatal("ACCESS_TOKEN_EXPIRE_TIME is not valid")
		}
	} else {
		log.Fatal("ACCESS_TOKEN_EXPIRE_TIME is not set")
	}

	if val, ok := os.LookupEnv("ACCESS_TOKEN_SECRET"); ok && len(val) > 0 {
		cfg.accessTokenSecret = val
	} else {
		log.Fatal("ACCESS_TOKEN_SECRET is not set")
	}

	if val, ok := os.LookupEnv("GOOGLE_OAUTH_CLIENT_ID"); ok && len(val) > 0 {
		cfg.googleOAuthClientId = val
	} else {
		log.Fatal("GOOGLE_OAUTH_CLIENT_ID is not set")
	}

	if val, ok := os.LookupEnv("GOOGLE_OAUTH_CLIENT_SECRET"); ok && len(val) > 0 {
		cfg.googleOAuthClientSecret = val
	} else {
		log.Fatal("GOOGLE_OAUTH_CLIENT_SECRET is not set")
	}

	if val, ok := os.LookupEnv("GOOGLE_OAUTH_REDIRECT_URL"); ok && len(val) > 0 {
		cfg.googleOAuthRedirectUrl = val
	} else {
		log.Fatal("GOOGLE_OAUTH_REDIRECT_URL is not set")
	}

	if val, ok := os.LookupEnv("GOOGLE_OAUTH_STATE"); ok && len(val) > 0 {
		cfg.googleOAuthState = val
	} else {
		log.Fatal("GOOGLE_OAUTH_STATE is not set")
	}

	if val, ok := os.LookupEnv("REDIS_MASTER_HOST"); ok && len(val) > 0 {
		cfg.redisMasterHost = val
	} else {
		log.Fatal("REDIS_MASTER_HOST is not set")
	}

	if val, ok := os.LookupEnv("REDIS_MASTER_PORT"); ok && len(val) > 0 {
		cfg.redisMasterPort = val
	} else {
		log.Fatal("REDIS_MASTER_PORT is not set")
	}

	if val, ok := os.LookupEnv("REDIS_MASTER_PASSWORD"); ok {
		cfg.redisMasterPassword = val
	}

	if val, ok := os.LookupEnv("REDIS_MASTER_DB"); ok && len(val) > 0 {
		db, err := strconv.Atoi(val)
		if err != nil {
			log.Fatal("REDIS_MASTER_DB is not valid")
		}

		cfg.redisMasterDb = db
	} else {
		log.Fatal("REDIS_MASTER_DB is not set")
	}

	if val, ok := os.LookupEnv("ASYNQ_REDIS_MASTER_HOST"); ok && len(val) > 0 {
		cfg.asynqRedisMasterHost = val
	} else {
		log.Fatal("ASYNQ_REDIS_MASTER_HOST is not set")
	}
	if val, ok := os.LookupEnv("ASYNQ_REDIS_MASTER_PORT"); ok && len(val) > 0 {
		cfg.asynqRedisMasterPort = val
	} else {
		log.Fatal("ASYNQ_REDIS_MASTER_PORT is not set")
	}

	if val, ok := os.LookupEnv("ASYNQ_REDIS_MASTER_PASSWORD"); ok {
		cfg.asynqRedisMasterPassword = val
	}

	if val, ok := os.LookupEnv("ASYNQ_REDIS_MASTER_DB"); ok && len(val) > 0 {
		db, err := strconv.Atoi(val)
		if err != nil {
			log.Fatal("ASYNQ_REDIS_MASTER_DB is invalid")
		}

		cfg.asynqRedisMasterDb = db
	} else {
		log.Fatal("ASYNC_REDIS_MASTER_DB is not set")
	}

	if val, ok := os.LookupEnv("SMTP_HOST"); ok && len(val) > 0 {
		cfg.smtpHost = val
	} else {
		log.Fatal("SMTP_HOST is not set")
	}

	if val, ok := os.LookupEnv("SMTP_PORT"); ok && len(val) > 0 {
		cfg.smtpPort = val
	} else {
		log.Fatal("SMTP_PORT is not set")
	}

	if val, ok := os.LookupEnv("SMTP_AUTH_USER"); ok && len(val) > 0 {
		cfg.smtpAuthUser = val
	} else {
		log.Fatal("SMTP_AUTH_USER is not set")
	}

	if val, ok := os.LookupEnv("SMTP_AUTH_PASSWORD"); ok && len(val) > 0 {
		cfg.smtpAuthPassword = val
	} else {
		log.Fatal("SMTP_AUTH_PASSWORD is not set")
	}

	if val, ok := os.LookupEnv("SMTP_SENDER_NAME"); ok && len(val) > 0 {
		cfg.smtpSenderName = val
	} else {
		log.Fatal("SMTP_SENDER_NAME is not set")
	}

	if val, ok := os.LookupEnv("WEB_BASE_URL"); ok && len(val) > 0 {
		cfg.webBaseUrl = val
	} else {
		log.Fatal("WEB_BASE_URL is not set")
	}

	if val, ok := os.LookupEnv("CENTRAL_WEB_BASE_URL"); ok && len(val) > 0 {
		cfg.centralWebBaseUrl = val
	} else {
		log.Fatal("CENTRAL_WEB_BASE_URL is not set")
	}

	return cfg
}
