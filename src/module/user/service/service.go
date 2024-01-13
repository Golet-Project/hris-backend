package service

import (
	"hroost/module/user/db"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	db *db.Db
}

type Dependency struct {
	Pg    *pgxpool.Pool
	Redis *redis.Client
}

func New(d *Dependency) *Service {
	if d.Pg == nil {
		log.Fatal("[x] Database connection required on user module")
	}
	if d.Redis == nil {
		log.Fatal("[x] Redis connection required on user module")
	}

	db := db.New(&db.Dependency{
		MasterConn: d.Pg,
		Redis:      d.Redis,
	})

	return &Service{
		db: db,
	}
}
