package repository

import (
	"context"

	"github.com/Aldiwildan77/backend-hexa-template/core/entity"
	"github.com/Aldiwildan77/backend-hexa-template/pkg/pagination"
)

type ReportRepository interface {
	GetReportByID(ctx context.Context, id string) (*entity.Report, error)
	GetReportsByQuery(ctx context.Context, param entity.GetReportsByQuery) ([]*entity.Report, *pagination.Pagination, error)
	Create(ctx context.Context, report *entity.Report) error
	Update(ctx context.Context, report *entity.Report) error
	Delete(ctx context.Context, id string) error
}
