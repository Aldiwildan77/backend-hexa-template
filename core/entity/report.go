package entity

import (
	"time"

	"github.com/Aldiwildan77/backend-hexa-template/pkg/pagination"
	"github.com/Aldiwildan77/backend-hexa-template/pkg/query"
)

type Report struct {
	ReportID    string     `json:"report_id" validate:"required"`
	ReporterID  int64      `json:"reporter_id"`
	ReportedID  int64      `json:"reported_id"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"-"`
	UpdatedAt   *time.Time `json:"-"`
	DeletedAt   *time.Time `json:"-"`
}

type GetReportsByQuery struct {
	Query string
	Sort  string
	Order string
	*pagination.Pagination
}

func (q *GetReportsByQuery) ValidateQuery() {
	if q.Sort == "" {
		q.Sort = "created_at"
	}

	if q.Order == "" {
		q.Order = query.OrderDesc
	}
}
