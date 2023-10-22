package internal

import (
	"hris/module/auth/internal/db"
	"hris/module/auth/internal/redis"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	redisClient "github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

var oauthState = os.Getenv("OAUTH_STATE")

type Internal struct {
	db    *db.Db
	redis *redis.Redis

	oauth2Cfg *oauth2.Config
}

type Dependency struct {
	Pg    *pgxpool.Pool
	Redis *redisClient.Client
}

func New(d *Dependency) *Internal {
	return &Internal{
		db: &db.Db{
			Pg:    d.Pg,
			Redis: d.Redis,
		},
		redis: &redis.Redis{
			Client: d.Redis,
		},

		oauth2Cfg: &oauth2.Config{
			ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
			ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
			Endpoint:     endpoints.Google,
			RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
		},
	}
}
