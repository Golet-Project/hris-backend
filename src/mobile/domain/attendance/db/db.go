package db

import (
	"context"
	"fmt"
	"hroost/infrastructure/store/postgres"
	"hroost/shared/primitive"
)

type IDbStore interface {
	AddAttendance(ctx context.Context, domain string, in AddAttendanceIn) error

	CheckTodayAttendanceById(ctx context.Context, domain, uid string, timezone primitive.Timezone) (exist bool, err error)

	CheckEmployeeById(ctx context.Context, domain, uid string) (exists bool, err error)

	Checkout(ctx context.Context, domain, uid string) (rowsAffected int64, err error)

	GetTodayAttendance(ctx context.Context, domain string, param GetTodayAttendanceIn) (out GetTodayAttendanceOut, err error)
}

type Db struct {
	pgResolver *postgres.Resolver
}

type Config struct {
	PgResolver *postgres.Resolver
}

func New(cfg *Config) (*Db, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config required")
	}
	if cfg.PgResolver == nil {
		return nil, fmt.Errorf("pgResolver required")
	}

	return &Db{
		pgResolver: cfg.PgResolver,
	}, nil
}
