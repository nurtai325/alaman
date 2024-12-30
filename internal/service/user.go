package service

import (
	"context"

	"github.com/nurtai325/alaman/internal/db/repository"
)

// TODO: service user model
// TODO: custom cells
// TODO: pagination complete
// TODO: generic functions for each resource
func (s *Service) GetUsers(ctx context.Context, minId, maxId int) ([]repository.User, error) {
	return s.queries.GetUsers(ctx, repository.GetUsersParams{
		ID:   int32(minId),
		ID_2: int32(maxId),
	})
}
