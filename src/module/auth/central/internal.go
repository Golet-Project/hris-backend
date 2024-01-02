package central

import (
	"hris/module/auth/central/db"
	"hris/module/auth/central/redis"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	redisClient "github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

var oauthState = os.Getenv("OAUTH_STATE")

type Central struct {
	db    *db.Db
	redis *redis.Redis

	oauth2Cfg *oauth2.Config
}

type Dependency struct {
	Pg    *pgxpool.Pool
	Redis *redisClient.Client
}

func New(d *Dependency) *Central {
	if d.Pg == nil {
		log.Fatal("[x] Database connection required on auth module")
	}
	if d.Redis == nil {
		log.Fatal("[x] Redis connection required on auth module")
	}

	dbImpl := db.New(&db.Dependency{
		MasterConn: d.Pg,
		Redis:      d.Redis,
	})
	redisImpl := redis.New(&redis.Dependency{
		Client: d.Redis,
	})

	return &Central{
		db:    dbImpl,
		redis: redisImpl,

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
