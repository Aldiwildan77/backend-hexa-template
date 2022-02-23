package report_repository

import (
	"time"

	"github.com/Aldiwildan77/backend-hexa-template/core/entity"
)

const (
	descriptionColumn = "description"
	reportIDColumn    = "report_id"
)

type Report struct {
	ReportID    string     `gorm:"column:report_id"`
	ReporterID  int64      `gorm:"column:reporter_id"`
	ReportedID  int64      `gorm:"column:reported_id"`
	Description string     `gorm:"column:description"`
	CreatedAt   *time.Time `gorm:"column:created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

func (r *Report) ToReportEntity() *entity.Report {
	if r == nil {
		return nil
	}

	return &entity.Report{
		ReportID:    r.ReportID,
		ReporterID:  r.ReporterID,
		ReportedID:  r.ReportedID,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		DeletedAt:   r.DeletedAt,
	}
}

func (Report) FromReportEntityToDTO(rep *entity.Report) *Report {
	return &Report{
		ReportID:    rep.ReportID,
		ReporterID:  rep.ReporterID,
		ReportedID:  rep.ReportedID,
		Description: rep.Description,
	}
}

func (Report) FromReportEntityToDTOWithCreatedAt(rep *entity.Report) *Report {
	return &Report{
		ReportID:    rep.ReportID,
		ReporterID:  rep.ReporterID,
		ReportedID:  rep.ReportedID,
		Description: rep.Description,
		CreatedAt:   rep.CreatedAt,
	}
}

func (Report) FromReportEntityToDTOWithCreatedAtAndUpdatedAt(rep *entity.Report) *Report {
	return &Report{
		ReportID:    rep.ReportID,
		ReporterID:  rep.ReporterID,
		ReportedID:  rep.ReportedID,
		Description: rep.Description,
		CreatedAt:   rep.CreatedAt,
		UpdatedAt:   rep.UpdatedAt,
	}
}

func (Report) TableName() string {
	return "report"
}
