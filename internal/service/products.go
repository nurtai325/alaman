package service

import (
	"context"
	"errors"
	"time"

	"github.com/nurtai325/alaman/internal/db/repository"
)

var (
	ErrInvalidProductCode = errors.New("штрих код тек сандардан тұруы керек")
)

type Product struct {
	Id         int
	Name       string
	InStock    int
	Price      int
	StockPrice int
	Code       string
	SaleCount  int
	CreatedAt  time.Time
}

func getSProduct(p repository.Product) Product {
	return Product{
		Id:         int(p.ID),
		Name:       p.Name,
		InStock:    int(p.InStock),
		Price:      int(p.Price),
		StockPrice: int(p.StockPrice),
		SaleCount:  int(p.SaleCount),
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

func (s *Service) GetProduct(ctx context.Context, id int) (Product, error) {
	p, err := s.queries.GetProduct(ctx, int32(id))
	if err != nil {
		return Product{}, err
	}
	return getSProduct(p), nil
}

func (s *Service) InsertProduct(ctx context.Context, name string, inStock, price, stockPrice, saleCount int) (Product, error) {
	p, err := s.queries.InsertProduct(ctx, repository.InsertProductParams{
		Name:       name,
		InStock:    int32(inStock),
		Price:      int32(price),
		StockPrice: int32(stockPrice),
		SaleCount:  int32(saleCount),
	})
	if err != nil {
		return Product{}, errors.Join(ErrInternal, err)
	}
	return getSProduct(p), nil
}

func (s *Service) UpdateProduct(ctx context.Context, name string, id, price, stockPrice, saleCount int) (Product, error) {
	p, err := s.queries.UpdateProduct(ctx, repository.UpdateProductParams{
		ID:         int32(id),
		Name:       name,
		Price:      int32(price),
		StockPrice: int32(stockPrice),
		SaleCount:  int32(saleCount),
	})
	if err != nil {
		return Product{}, errors.Join(ErrInternal, err)
	}
	return getSProduct(p), nil
}

func (s *Service) DeleteProduct(ctx context.Context, id int) (Product, error) {
	if err := validId(id); err != nil {
		return Product{}, err
	}
	p, err := s.queries.DeleteProduct(ctx, int32(id))
	if err != nil {
		return Product{}, err
	}
	return getSProduct(p), nil
}

func (s *Service) AddStockProduct(ctx context.Context, id, quantity int) (int, error) {
	inStock, err := s.queries.AddStockProduct(ctx, repository.AddStockProductParams{
		ID:      int32(id),
		InStock: int32(quantity),
	})
	if err != nil {
		return 0, err
	}
	_, err = s.queries.InsertProductChange(ctx, repository.InsertProductChangeParams{
		Quantity:  int32(quantity),
		IsIncome:  true,
		ProductID: int32(id),
	})
	if err != nil {
		return 0, err
	}
	return int(inStock), err
}

func (s *Service) RemoveStockProduct(ctx context.Context, id, quantity int) (int, error) {
	inStock, err := s.queries.RemoveStockProduct(ctx, repository.RemoveStockProductParams{
		ID:      int32(id),
		InStock: int32(quantity),
	})
	if err != nil {
		return 0, err
	}
	_, err = s.queries.InsertProductChange(ctx, repository.InsertProductChangeParams{
		Quantity:  int32(quantity),
		IsIncome:  false,
		ProductID: int32(id),
	})
	if err != nil {
		return 0, err
	}
	return int(inStock), err
}

func validProductCode(code string) bool {
	for _, r := range code {
		if r < 48 || r > 57 {
			return false
		}
	}
	return true
}
