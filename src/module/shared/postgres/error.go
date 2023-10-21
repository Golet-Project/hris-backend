package postgres

import "errors"

var ErrNoConnection = errors.New("connection doesn't exists")
var ErrConnectionAlreadyExists = errors.New("connection already exists")