package services

import (
	// "github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type contextKey int

const (
	// AcceptTypeKey is the context key used to store the value of the HTTP
	// Accept-Type header
	AcceptTypeKey contextKey = iota + 1
	// AuthorizationKey is the context key used to store the value of the HTTP
	// Authorization Header
	AuthorizationKey
)

// Env ...
type Env struct {
	// DB  *sqlx.DB
	Log *zerolog.Logger
}
