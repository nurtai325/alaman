package service

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nurtai325/alaman/internal/db/repository"
)

type Sale struct {
	Id           int
	UserId       int
	UserName     string
	FullPrice    float32
	DeliveryCost float32
	LoanCost     float32
	PaymentAt    time.Time
	CreatedAt    time.Time
}

func getSSale(sale repository.Sale) Sale {
	return Sale{
		Id:           int(sale.ID),
		FullPrice:    sale.FullSum,
		DeliveryCost: sale.DeliveryCost,
		LoanCost:     sale.LoanCost,
		PaymentAt:    sale.PaymentAt.Time,
		CreatedAt:    sale.CreatedAt.Time,
	}
}

func (s *Service) GetWeekSum(ctx context.Context) (int, error) {
	now, err := time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
	if err != nil {
		return 0, err
	}
	weekDay := now.Weekday()
	difference := weekDay - 1
	if weekDay == time.Sunday {
		difference = 6
	}
	weekStart := now.AddDate(0, 0, -int(difference))
	sum, err := s.queries.GetSum(ctx, pgtype.Timestamptz{
		Time:  weekStart,
		Valid: true,
	})
	if err != nil {
		return 0, err
	}
	return int(sum), nil
}

func (s *Service) GetMonthSum(ctx context.Context) (int, error) {
	now, err := time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
	if err != nil {
		return 0, err
	}
	monthStart := now.AddDate(0, 0, -now.Day())
	sum, err := s.queries.GetSum(ctx, pgtype.Timestamptz{
		Time:  monthStart,
		Valid: true,
	})
	if err != nil {
		return 0, err
	}
	return int(sum), nil
}

func (s *Service) GetScr(ctx context.Context) (float32, error) {
	now, err := time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
	if err != nil {
		return 0, err
	}
	weekDay := now.Weekday()
	difference := weekDay - 1
	if weekDay == time.Sunday {
		difference = 6
	}
	weekStart := now.AddDate(0, 0, -int(difference))
	leadCount, err := s.queries.GetLeadCount(ctx, pgtype.Timestamptz{
		Time:  weekStart,
		Valid: true,
	})
	if err != nil {
		return 0, err
	}
	soldLeadCount, err := s.queries.GetSoldLeadCount(ctx, pgtype.Timestamptz{
		Time:  weekStart,
		Valid: true,
	})
	if err != nil {
		return 0, err
	}
	if leadCount == 0 || soldLeadCount == 0 {
		return 0, nil
	}
	if leadCount == 0 {
		return 0, nil
	}
	return float32(soldLeadCount) / float32(leadCount) * 100, nil
}

func (s *Service) GetNewLeadsCount(ctx context.Context, weekStart time.Time) (int, error) {
	count, err := s.queries.GetNewLeadsCount(ctx, pgtype.Timestamptz{
		Time:  weekStart,
		Valid: true,
	})
	return int(count), err
}

type ChartsData struct {
	Week       []barData
	Month      []barData
	Manager    []barData
	Product    []barData
	WeekSum    int
	MonthSum   int
	Scr        float32
	AverageSum int
	NewLeads   int
	SalesCount int
}

type barData struct {
	Id     int
	Label  string
	Amount int
}

func (s *Service) GetSalesData(ctx context.Context) (*ChartsData, error) {
	scr, err := s.GetScr(ctx)
	if err != nil {
		return nil, err
	}
	now, err := time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
	if err != nil {
		return nil, err
	}
	weekDay := now.Weekday()
	difference := weekDay - 1
	if weekDay == time.Sunday {
		difference = 6
	}
	weekStart := now.AddDate(0, 0, -int(difference))
	monthStart := now.AddDate(0, 0, -now.Day())
	var start time.Time
	if weekStart.Before(monthStart) {
		start = weekStart
	} else {
		start = monthStart
	}
	chartsData := &ChartsData{}
	chartsData.Month = make([]barData, 31)
	chartsData.Week = []barData{
		{Label: "Дүйсенбі"},
		{Label: "Сейсенбі"},
		{Label: "Сәрсенбі"},
		{Label: "Бейсенбі"},
		{Label: "Жұма"},
		{Label: "Сенбі"},
		{Label: "Жексенбі"},
	}
	chartsData.Manager = make([]barData, 0, 10)
	chartsData.Product = make([]barData, 0)
	chartsData.Scr = scr
	newLeadsCount, err := s.GetNewLeadsCount(ctx, weekStart)
	if err != nil {
		return nil, err
	}
	chartsData.NewLeads = newLeadsCount
	sales, err := s.queries.GetSales(ctx, pgtype.Timestamptz{
		Time:  start,
		Valid: true,
	})
	if err != nil {
		return nil, err
	}
	for i := range chartsData.Month {
		chartsData.Month[i].Label = strconv.Itoa(i + 1)
	}
	checkCount := 0
	for _, sale := range sales {
		if sale.PaymentAt.Time.After(weekStart) {
			checkCount++
			day := sale.PaymentAt.Time.Weekday() - 1
			if sale.PaymentAt.Time.Weekday() == time.Sunday {
				day = 6
			}
			chartsData.Week[day].Amount += int(sale.ItemsSum)
			chartsData.WeekSum += int(sale.ItemsSum)
		}
		if !sale.PaymentAt.Time.After(monthStart) {
			continue
		}
		chartsData.Month[sale.PaymentAt.Time.Day()-1].Amount += int(sale.ItemsSum)
		chartsData.MonthSum += int(sale.ItemsSum)
		userFound := false
		for i, manager := range chartsData.Manager {
			if int32(manager.Id) == sale.UserID {
				userFound = true
				chartsData.Manager[i].Amount += int(sale.ItemsSum)
				break
			}
		}
		if !userFound {
			chartsData.Manager = append(chartsData.Manager, barData{
				Id:     int(sale.UserID),
				Label:  sale.UserName,
				Amount: int(sale.ItemsSum),
			})
		}
	}
	if checkCount == 0 {
		chartsData.AverageSum = 0
	} else {
		chartsData.AverageSum = chartsData.WeekSum / checkCount
	}
	salesItems, err := s.queries.GetSaleItemsByTime(ctx, pgtype.Timestamptz{
		Time:  weekStart,
		Valid: true,
	})
	if err != nil {
		return nil, err
	}
	salesCount := 0
	for _, saleItem := range salesItems {
		salesCount += int(saleItem.SaleCount)
		productFound := false
		for i, product := range chartsData.Product {
			if int32(product.Id) == saleItem.ProductID {
				productFound = true
				chartsData.Product[i].Amount += int(saleItem.Quantity * int32(saleItem.Price))
				break
			}
		}
		if !productFound {
			chartsData.Product = append(chartsData.Product, barData{
				Id:     int(saleItem.ProductID),
				Label:  saleItem.ProductName,
				Amount: int(saleItem.Quantity * int32(saleItem.Price)),
			})
		}
	}
	chartsData.SalesCount = salesCount
	return chartsData, nil
}
