package service

import (
	"errors"

	"github.com/nurtai325/alaman/internal/db/repository"
)

var (
	ErrFailedFormValidation = errors.New("failed form validation")
)

type Service struct {
	queries *repository.Queries
}

func New(queries *repository.Queries) *Service {
	return &Service{queries: queries}
}
