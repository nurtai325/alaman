package service

import (
	"errors"

	"github.com/nurtai325/alaman/internal/db/repository"
)

var (
	ErrInternal      = errors.New("internal err")
	ErrInvalidId     = errors.New("id is equal or less than zero")
	ErrInvalidArgs   = errors.New("invalid args")
	ErrInvalidOffset = errors.New("offset is less than 0")
	ErrInvalidLimit  = errors.New("limit is less than or equal to 0")
)

type Service struct {
	queries *repository.Queries
}

func New(queries *repository.Queries) *Service {
	return &Service{queries: queries}
}
