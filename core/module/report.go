package module

import (
	"context"

	"github.com/Aldiwildan77/backend-hexa-template/core/entity"
	"github.com/Aldiwildan77/backend-hexa-template/core/repository"
	"github.com/Aldiwildan77/backend-hexa-template/pkg/pagination"
)

type ReportUsecase interface {
	CreateReport(ctx context.Context, report *entity.Report) (*entity.Report, error)
	UpdateReport(ctx context.Context, report *entity.Report) (*entity.Report, error)
	DeleteReport(ctx context.Context, id string) error
	GetReport(ctx context.Context, id string) (*entity.Report, error)
	GetReports(ctx context.Context, param entity.GetReportsByQuery) ([]*entity.Report, *pagination.Pagination, error)
}

type reportUsecase struct {
	reportRepository repository.ReportRepository
}

func NewReportUsecase(reportRepository repository.ReportRepository) ReportUsecase {
	return &reportUsecase{
		reportRepository: reportRepository,
	}
}

func (uc *reportUsecase) CreateReport(ctx context.Context, report *entity.Report) (*entity.Report, error) {
	if report.ReportID == "" {
		return nil, entity.ErrReportIDRequired
	}

	if report.ReportedID == 0 {
		return nil, entity.ErrReportedIDRequired
	}

	if report.ReporterID == 0 {
		return nil, entity.ErrReporterIDRequired
	}

	return report, uc.reportRepository.Create(ctx, report)
}

func (uc *reportUsecase) UpdateReport(ctx context.Context, report *entity.Report) (*entity.Report, error) {
	if report.ReportID == "" {
		return nil, entity.ErrReportIDRequired
	}

	if report.ReportedID == 0 {
		return nil, entity.ErrReportedIDRequired
	}

	if report.ReporterID == 0 {
		return nil, entity.ErrReporterIDRequired
	}

	_, err := uc.reportRepository.GetReportByID(ctx, report.ReportID)
	if err != nil {
		return nil, err
	}

	return report, uc.reportRepository.Update(ctx, report)
}

func (uc *reportUsecase) DeleteReport(ctx context.Context, id string) error {
	if id == "" {
		return entity.ErrReportIDRequired
	}

	_, err := uc.reportRepository.GetReportByID(ctx, id)
	if err != nil {
		return err
	}

	return uc.reportRepository.Delete(ctx, id)
}

func (uc *reportUsecase) GetReport(ctx context.Context, id string) (*entity.Report, error) {
	if id == "" {
		return nil, entity.ErrReportIDRequired
	}

	return uc.reportRepository.GetReportByID(ctx, id)
}

func (uc *reportUsecase) GetReports(ctx context.Context, param entity.GetReportsByQuery) ([]*entity.Report, *pagination.Pagination, error) {
	return uc.reportRepository.GetReportsByQuery(ctx, param)
}
