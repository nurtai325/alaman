package service

import (
	"context"
	"time"

	"github.com/nurtai325/alaman/internal/db/repository"
)

type Product struct {
	Id         int
	Name       string
	InStock    int
	Price      int
	StockPrice int
	CreatedAt  time.Time
}

func getSProduct(p repository.Product) Product {
	return Product{
		Id:         int(p.ID),
		Name:       p.Name,
		InStock:    int(p.InStock),
		Price:      int(p.Price),
		StockPrice: int(p.StockPrice),
		CreatedAt:  p.CreatedAt.Time,
	}
}

func (s *Service) GetProducts(ctx context.Context, offset, limit int) ([]Product, error) {
	if offset < 0 {
		return nil, ErrInvalidOffset
	} else if limit <= 0 {
		return nil, ErrInvalidLimit
	}
	products, err := s.queries.GetProducts(ctx, repository.GetProductsParams{
		Offset: int64(offset),
		Limit:  int64(limit),
	})
	if err != nil {
		return nil, err
	}
	sProducts := make([]Product, 0, len(products))
	for _, product := range products {
		sProducts = append(sProducts, getSProduct(product))
	}
	return sProducts, nil
}
