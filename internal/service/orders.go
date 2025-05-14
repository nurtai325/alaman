package service

import (
	"context"
	"errors"
	"fmt"
)

type FullLead struct {
}

func (s *Service) GetOrder(ctx context.Context, id int) (Lead, error) {
	fullLead, err := s.queries.GetFullLead(ctx, int32(id))
	if err != nil {
		return Lead{}, errors.Join(ErrInternal, err)
	}
	items, err := s.queries.GetSaleItems(ctx, fullLead.SaleID.Int32)
	sItems := make([]SaleItem, 0, len(items))
	itemsStr := ""
	for _, item := range items {
		sItems = append(sItems, SaleItem{
			Id:          int(item.ID),
			ProductName: item.ProductName,
			Quantity:    int(item.Quantity),
			ProductId:   int(item.ProductID),
		})
		itemsStr += fmt.Sprintf("%d %s\n", item.Quantity, item.ProductName)
	}
	lead := Lead{
		Id:                 int(fullLead.ID),
		Phone:              fullLead.Phone,
		Address:            fullLead.Address.String,
		Name:               fullLead.Name.String,
		UserName:           fullLead.UserName,
		FullPrice:          fullLead.FullSum,
		DeliveryCost:       fullLead.DeliveryCost,
		LoanCost:           fullLead.LoanCost,
		SaleType:           saleType(fullLead.SaleType),
		SaleTypeName:       getSaleTypeName(saleType(fullLead.SaleType)),
		Items:              sItems,
		PaymentAt:          fullLead.SoldAt.Time,
		PaymentAtFormatted: fullLead.PaymentAt.Time.Format(dateTimeFormat),
		DeliveryType:       deliveryType(fullLead.DeliveryType.String),
		DeliveryTypeName:   getDeliveryTypeName(deliveryType(fullLead.DeliveryType.String)),
	}
	return lead, nil
}
