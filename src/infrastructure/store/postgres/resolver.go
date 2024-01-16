package postgres

import (
	"errors"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNoConnection = errors.New("connection doesn't exists")
var ErrConnectionAlreadyExists = errors.New("connection already exists")

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

	conn sync.Map
}

func NewResolver() *Resolver {
	return &Resolver{
		conn: sync.Map{},
	}
}

func (r *Resolver) Register(d Database) error {
	_, ok := r.conn.Load(d.DomainName)
	if ok {
		return ErrConnectionAlreadyExists
	}

	r.conn.Store(d.DomainName, d.Pool)
	return nil
}

func (r *Resolver) Resolve(domain Domain) (*pgxpool.Pool, error) {
	conn, exists := r.conn.Load(domain)
	if !exists {
		return nil, ErrNoConnection
	}

	return conn.(*pgxpool.Pool), nil
}

func (r *Resolver) Remove(domain Domain) error {
	_, ok := r.conn.Load(domain)
	if !ok {
		return ErrNoConnection
	}

	r.conn.Delete(domain)

	return nil
}

func (r *Resolver) GetAllConn() *sync.Map {
	return &r.conn
}
