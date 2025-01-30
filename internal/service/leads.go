package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nurtai325/alaman/internal/config"
	"github.com/nurtai325/alaman/internal/db"
	"github.com/nurtai325/alaman/internal/db/repository"
)

var (
	ErrInvalidSaleType = errors.New("invalid sale type")
)

type Lead struct {
	Id           int
	Name         string
	UserName     string
	Address      string
	Phone        string
	Completed    bool
	SaleType     saleType
	FullPrice    float32
	DeliveryCost float32
	LoanCost     float32
	Items        []SaleItem
	UserId       int
	SaleId       int
	CreatedAt    time.Time
}

type saleType string

const (
	kaspiLoan     saleType = "kaspi-loan"
	cash          saleType = "cash"
	kaspiRed      saleType = "red"
	kaspiTransfer saleType = "kaspi-transfer"
	kaspiQr       saleType = "kaspi-qr"
)

func getSLead(lead repository.Lead) Lead {
	return Lead{
		Id:        int(lead.ID),
		Name:      lead.Name.String,
		Address:   lead.Address.String,
		Phone:     lead.Phone,
		Completed: lead.Completed,
		UserId:    int(lead.UserID.Int32),
		SaleId:    int(lead.SaleID.Int32),
		CreatedAt: lead.CreatedAt.Time,
	}
}

func (s *Service) GetNewLeads(ctx context.Context) ([]Lead, error) {
	leads, err := s.queries.GetNewLeads(ctx)
	if err != nil {
		return nil, err
	}
	sLeads := make([]Lead, 0, len(leads))
	for _, lead := range leads {
		sLeads = append(sLeads, getSLead(lead))
	}
	return sLeads, nil
}

func (s *Service) GetAssignedLeads(ctx context.Context) ([]Lead, error) {
	leads, err := s.queries.GetAssignedLeads(ctx)
	if err != nil {
		return nil, err
	}
	sLeads := make([]Lead, 0, len(leads))
	for _, lead := range leads {
		sLeads = append(sLeads, Lead{
			Id:        int(lead.ID),
			Name:      lead.Name.String,
			UserName:  lead.UserName,
			Address:   lead.Address.String,
			Phone:     lead.Phone,
			Completed: lead.Completed,
			UserId:    int(lead.UserID.Int32),
			SaleId:    int(lead.SaleID.Int32),
			CreatedAt: lead.CreatedAt.Time,
		})
	}
	return sLeads, nil
}

func (s *Service) GetInDeliveryLeads(ctx context.Context) ([]Lead, error) {
	leads, err := s.queries.GetInDeliveryLeads(ctx)
	if err != nil {
		return nil, err
	}
	sLeads := make([]Lead, 0, len(leads))
	for _, lead := range leads {
		items, err := s.queries.GetSaleItems(ctx, lead.SaleID.Int32)
		if err != nil {
			return nil, err
		}
		sItems := make([]SaleItem, 0, len(items))
		for _, item := range items {
			sItems = append(sItems, SaleItem{
				Id:          int(item.ID),
				ProductName: item.ProductName,
				ProductId:   int(item.ProductID),
				Quantity:    int(item.Quantity),
			})
		}
		sLeads = append(sLeads, Lead{
			Id:        int(lead.ID),
			Name:      lead.Name.String,
			UserName:  lead.UserName,
			Address:   lead.Address.String,
			Phone:     lead.Phone,
			Completed: lead.Completed,
			UserId:    int(lead.UserID.Int32),
			SaleId:    int(lead.SaleID.Int32),
			Items:     sItems,
			CreatedAt: lead.CreatedAt.Time,
		})
	}
	return sLeads, nil
}

func (s *Service) GetCompletedLeads(ctx context.Context) ([]Lead, error) {
	leads, err := s.queries.GetCompletedLeads(ctx)
	if err != nil {
		return nil, err
	}
	sLeads := make([]Lead, 0, len(leads))
	for _, lead := range leads {
		items, err := s.queries.GetSaleItems(ctx, lead.SaleID.Int32)
		if err != nil {
			return nil, err
		}
		sItems := make([]SaleItem, 0, len(items))
		for _, item := range items {
			sItems = append(sItems, SaleItem{
				Id:          int(item.ID),
				ProductName: item.ProductName,
				ProductId:   int(item.ProductID),
				Quantity:    int(item.Quantity),
			})
		}
		sLeads = append(sLeads, Lead{
			Id:        int(lead.ID),
			Name:      lead.Name.String,
			UserName:  lead.UserName,
			Address:   lead.Address.String,
			Phone:     lead.Phone,
			Completed: lead.Completed,
			UserId:    int(lead.UserID.Int32),
			SaleId:    int(lead.SaleID.Int32),
			Items:     sItems,
			CreatedAt: lead.CreatedAt.Time,
		})
	}
	return sLeads, nil
}

func (s *Service) InsertLead(ctx context.Context, phone string) (Lead, error) {
	if !validPhone(phone) {
		return Lead{}, ErrInvalidPhone
	}
	lead, err := s.queries.InsertLead(ctx, phone)
	if err != nil {
		return Lead{}, err
	}
	return getSLead(lead), nil
}

func (s *Service) AssignLead(ctx context.Context, id, userId int) (Lead, error) {
	lead, err := s.queries.AssignLead(ctx, repository.AssignLeadParams{
		ID: int32(id),
		UserID: pgtype.Int4{
			Int32: int32(userId),
			Valid: true,
		},
	})
	if err != nil {
		return Lead{}, err
	}
	return getSLead(lead), nil
}

func (s *Service) CompleteLead(ctx context.Context, id int) error {
	_, err := s.queries.CompleteLead(ctx, int32(id))
	return err
}

type SaleItem struct {
	Id          int
	ProductName string
	ProductId   int
	Quantity    int
}

type SellLeadParams struct {
	Id           int
	Name         string
	Address      string
	Type         string
	FullSum      float32
	DeliveryCost float32
	LoanCost     float32
	ItemsSum     float32
	Items        []SaleItem
}

func (s *Service) SellLead(ctx context.Context, arg SellLeadParams) (Lead, error) {
	if !validSaleType(arg.Type) {
		return Lead{}, fmt.Errorf("%w: %s", ErrInvalidSaleType, arg.Type)
	}
	conf, err := config.New()
	if err != nil {
		return Lead{}, errors.Join(ErrInternal, err)
	}
	pool, err := db.New(conf)
	if err != nil {
		return Lead{}, errors.Join(ErrInternal, err)
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return Lead{}, errors.Join(ErrInternal, err)
	}
	defer tx.Rollback(ctx)
	q := s.queries.WithTx(tx)
	_, err = q.SetLeadInfo(ctx, repository.SetLeadInfoParams{
		ID: int32(arg.Id),
		Name: pgtype.Text{
			String: arg.Name,
			Valid:  true,
		},
		Address: pgtype.Text{
			String: arg.Address,
			Valid:  true,
		},
	})
	if err != nil {
		return Lead{}, errors.Join(ErrInternal, err)
	}
	sale, err := q.InsertSale(ctx, repository.InsertSaleParams{
		Type:         string(arg.Type),
		FullSum:      arg.FullSum,
		DeliveryCost: arg.DeliveryCost,
		LoanCost:     arg.LoanCost,
		ItemsSum:     arg.ItemsSum,
	})
	if err != nil {
		return Lead{}, errors.Join(ErrInternal, err)
	}
	for _, item := range arg.Items {
		_, err := q.InsertSaleItem(ctx, repository.InsertSaleItemParams{
			SaleID:    sale.ID,
			ProductID: int32(item.ProductId),
			Quantity:  int32(item.Quantity),
		})
		if err != nil {
			return Lead{}, errors.Join(ErrInternal, err)
		}
		_, err = q.RemoveStockProduct(ctx, repository.RemoveStockProductParams{
			ID:      int32(item.ProductId),
			InStock: int32(item.Quantity),
		})
		if err != nil {
			return Lead{}, errors.Join(ErrInternal, err)
		}
	}
	lead, err := q.SellLead(ctx, repository.SellLeadParams{
		ID: int32(arg.Id),
		SaleID: pgtype.Int4{
			Int32: sale.ID,
			Valid: true,
		},
	})
	if err != nil {
		return Lead{}, errors.Join(ErrInternal, err)
	}
	fullLead, err := q.GetFullLead(ctx, lead.ID)
	if err != nil {
		return Lead{}, errors.Join(ErrInternal, err)
	}
	items, err := q.GetSaleItems(ctx, fullLead.SaleID.Int32)
	sItems := make([]SaleItem, 0, len(items))
	for _, item := range items {
		sItems = append(sItems, SaleItem{
			Id:          int(item.ID),
			ProductName: item.ProductName,
			ProductId:   int(item.ProductID),
			Quantity:    int(item.Quantity),
		})
	}
	return Lead{
		Id:       int(fullLead.ID),
		Phone:    fullLead.Phone,
		Address:  fullLead.Address.String,
		Name:     fullLead.Name.String,
		UserName: fullLead.UserName,
		Items:    sItems,
	}, tx.Commit(ctx)
}

func validSaleType(saletypeStr string) bool {
	saletype := saleType(saletypeStr)
	switch saletype {
	case kaspiLoan:
		return true
	case cash:
		return true
	case kaspiRed:
		return true
	case kaspiTransfer:
		return true
	case kaspiQr:
		return true
	default:
		return false
	}
}
