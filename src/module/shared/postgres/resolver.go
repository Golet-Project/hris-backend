package postgres

import (
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Domain string

const (
	MasterDomain Domain = "master"
)

type Database struct {
	DomainName Domain
	Pool       *pgxpool.Pool
}

type Resolver struct {
	MasterConn *pgxpool.Pool

	tenantConn sync.Map
}

func NewResolver(masterConn *pgxpool.Pool) *Resolver {
	return &Resolver{
		MasterConn: masterConn,
		tenantConn: sync.Map{},
	}
}

func (r *Resolver) Register(d Database) error {
	_, ok := r.tenantConn.Load(d.DomainName)
	if ok {
		return ErrConnectionAlreadyExists
	}

	r.tenantConn.Store(d.DomainName, d.Pool)
	return nil
}

func (r *Resolver) Resolve(domain Domain) (*pgxpool.Pool, error) {
	conn, exists := r.tenantConn.Load(domain)
	if !exists {
		return nil, ErrNoConnection
	}

	return conn.(*pgxpool.Pool), nil
}

func (r *Resolver) Remove(domain Domain) error {
	_, ok := r.tenantConn.Load(domain)
	if !ok {
		return ErrNoConnection
	}

	r.tenantConn.Delete(domain)

	return nil
}

func (r *Resolver) GetAllTenantConn() *sync.Map {
	return &r.tenantConn
}
