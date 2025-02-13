package service

import (
	"errors"

	"github.com/nurtai325/alaman/internal/db/repository"
)

var (
	ErrInternal      = errors.New("internal err")
	ErrNotFound      = errors.New("resource not found")
	ErrInvalidId     = errors.New("id is equal or less than zero")
	ErrInvalidArgs   = errors.New("invalid args")
	ErrInvalidOffset = errors.New("offset is less than 0")
	ErrInvalidLimit  = errors.New("limit is less than or equal to 0")
	ErrAlreadySold   = errors.New("lead is already sold")
)

const (
	dateTimeFormat = "2006/01/02 15:04"
	dateFormat     = "2006/01/02"
)

type Service struct {
	queries *repository.Queries
}

func New(queries *repository.Queries) *Service {
	return &Service{queries: queries}
}
