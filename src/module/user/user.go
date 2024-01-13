package user

import (
	"hroost/module/user/service"
	"log"

	redisClient "github.com/redis/go-redis/v9"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	UserService *service.Service
}

type Dependency struct {
	MasterDB    *pgxpool.Pool
	RedisClient *redisClient.Client
}

func InitUser(d *Dependency) *User {
	if d.MasterDB == nil {
		log.Fatal("[x] User package require a database connection")
	}
	if d.RedisClient == nil {
		log.Fatal("[x] User package require a redis connection")
	}

	userService := service.New(&service.Dependency{
		Pg:    d.MasterDB,
		Redis: d.RedisClient,
	})

	return &User{
		UserService: userService,
	}
}
