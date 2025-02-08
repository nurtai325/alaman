package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nurtai325/alaman/internal/db/repository"
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

func (s *Service) InsertReport(ctx context.Context, name string, startAt, endAt time.Time) (Report, error) {
	report, err := s.queries.InsertReport(ctx, repository.InsertReportParams{
		Name: name,
		Path: "",
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
