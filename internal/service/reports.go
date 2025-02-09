package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nurtai325/alaman/internal/db/repository"
	"github.com/xuri/excelize/v2"
)

type Report struct {
	Id           int
	Name         string
	Path         string
	StartAt      time.Time
	EndAt        time.Time
	Span         string
	CreatedAt    time.Time
	CreatedAtStr string
}

func timeSpan(s, e time.Time) string {
	return fmt.Sprintf("%s-%s", s.Format(dateFormat), e.Format(dateFormat))
}

func getSReport(r repository.Report) Report {
	return Report{
		Id:           int(r.ID),
		Name:         r.Name,
		Path:         r.Path,
		StartAt:      r.StartAt.Time,
		EndAt:        r.EndAt.Time,
		Span:         timeSpan(r.StartAt.Time, r.EndAt.Time),
		CreatedAt:    r.CreatedAt.Time,
		CreatedAtStr: r.CreatedAt.Time.Format(dateTimeFormat),
	}
}

func (s *Service) GetReports(ctx context.Context) ([]Report, error) {
	reports, err := s.queries.GetReports(ctx)
	if err != nil {
		return nil, err
	}
	sReports := make([]Report, 0, len(reports))
	for _, r := range reports {
		sReports = append(sReports, getSReport(r))
	}
	return sReports, nil
}

func (s *Service) GetReport(ctx context.Context, id int) (Report, error) {
	report, err := s.queries.GetReport(ctx, int32(id))
	if err != nil {
		return Report{}, err
	}
	return getSReport(report), nil

}

const (
	reportsDir = "./assets/reports"
)

type productReport struct {
	Order     int
	Incoming  int
	Outcoming int
	Name      string
	SaleCount int
	Sold      int
	InStock   int
	SoldSum   float32
}

func (s *Service) InsertReport(ctx context.Context, name string, startAt, endAt time.Time) (Report, error) {
	products, err := s.queries.GetProducts(ctx, repository.GetProductsParams{
		Offset: 0,
		Limit:  1000,
	})
	productReports := make([]productReport, 0, len(products))
	totalSold := 0
	totalSaleCount := 0
	totalSoldSum := 0
	for i, product := range products {
		startAt = startAt.Add(time.Hour * 5)
		endAt = endAt.Add(time.Hour * 24)
		productSaleSum, err := s.queries.GetReportByProduct(ctx, repository.GetReportByProductParams{
			ProductID: product.ID,
			CreatedAt: pgtype.Timestamptz{
				Time:  startAt,
				Valid: true,
			},
			CreatedAt_2: pgtype.Timestamptz{
				Time:  endAt,
				Valid: true,
			},
		})
		if err != nil {
			return Report{}, err
		}
		productIncoming, err := s.queries.GetProductIncoming(ctx, repository.GetProductIncomingParams{
			ProductID: product.ID,
			CreatedAt: pgtype.Timestamptz{
				Time:  startAt,
				Valid: true,
			},
			CreatedAt_2: pgtype.Timestamptz{
				Time:  endAt,
				Valid: true,
			},
		})
		productOutcoming, err := s.queries.GetProductOutcoming(ctx, repository.GetProductOutcomingParams{
			ProductID: product.ID,
			CreatedAt: pgtype.Timestamptz{
				Time:  startAt,
				Valid: true,
			},
			CreatedAt_2: pgtype.Timestamptz{
				Time:  endAt,
				Valid: true,
			},
		})
		totalSold += int(productSaleSum.Sold.Int64)
		totalSaleCount += int(productSaleSum.SaleCountSum.Int64)
		totalSoldSum += int(productSaleSum.SoldSum.Float32)
		productReports = append(productReports, productReport{
			Order:     i + 1,
			Incoming:  int(productIncoming),
			Outcoming: int(productOutcoming),
			Name:      product.Name,
			SaleCount: int(productSaleSum.SaleCountSum.Int64),
			Sold:      int(productSaleSum.Sold.Int64),
			InStock:   int(product.InStock),
			SoldSum:   productSaleSum.SoldSum.Float32,
		})
	}
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()
	sheet1 := "Sheet1"
	firstRow := &[]any{name, "", startAt.Format(time.DateOnly), "-", endAt.Format(time.DateOnly)}
	err = f.SetSheetRow(sheet1, "A1", firstRow)
	if err != nil {
		return Report{}, err
	}
	reportHeaders := &[]any{"#", "Аты", "Келді", "Кетті", "Қалды", "Сатылды", "Сатылым Саны", "Сатылған сумма"}
	err = f.SetSheetRow(sheet1, "A3", reportHeaders)
	if err != nil {
		return Report{}, err
	}
	currentRow := 4
	for _, productReport := range productReports {
		row := &[]any{productReport.Order, productReport.Name, productReport.Incoming, productReport.Outcoming, productReport.InStock, productReport.Sold, productReport.SaleCount, productReport.SoldSum}
		err = f.SetSheetRow(sheet1, fmt.Sprintf("A%d", currentRow), row)
		if err != nil {
			return Report{}, err
		}
		currentRow += 1
	}
	err = f.SetCellInt(sheet1, fmt.Sprintf("F%d", currentRow), totalSold)
	if err != nil {
		return Report{}, err
	}
	err = f.SetCellInt(sheet1, fmt.Sprintf("G%d", currentRow), totalSaleCount)
	if err != nil {
		return Report{}, err
	}
	err = f.SetCellInt(sheet1, fmt.Sprintf("H%d", currentRow), totalSoldSum)
	if err != nil {
		return Report{}, err
	}
	fPath := fmt.Sprintf("%s/%s-%s-%d.xlsx", reportsDir, startAt.Format(time.DateOnly), endAt.Format(time.DateOnly), time.Now().UnixNano())
	err = f.SaveAs(fPath)
	if err != nil {
		return Report{}, err
	}
	report, err := s.queries.InsertReport(ctx, repository.InsertReportParams{
		Name: name,
		Path: fPath,
		StartAt: pgtype.Timestamptz{
			Time:  startAt,
			Valid: true,
		},
		EndAt: pgtype.Timestamptz{
			Time:  endAt,
			Valid: true,
		},
	})
	if err != nil {
		return Report{}, err
	}
	return getSReport(report), nil

}

func (s *Service) UpdateReport(ctx context.Context, id int, name string) (Report, error) {
	report, err := s.queries.UpdateReport(ctx, repository.UpdateReportParams{
		ID:   int32(id),
		Name: name,
	})
	if err != nil {
		return Report{}, err
	}
	return getSReport(report), nil
}

func (s *Service) DeleteReport(ctx context.Context, id int) (Report, error) {
	report, err := s.queries.DeleteReport(ctx, int32(id))
	if err != nil {
		return Report{}, err
	}
	return getSReport(report), nil
}
