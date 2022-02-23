package report_repository

import (
	"context"
	"fmt"

	"github.com/Aldiwildan77/backend-hexa-template/core/entity"
	rp "github.com/Aldiwildan77/backend-hexa-template/core/repository"
	"github.com/Aldiwildan77/backend-hexa-template/pkg/pagination"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) rp.ReportRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetReportByID(ctx context.Context, id string) (*entity.Report, error) {
	q := fmt.Sprintf("%s = '%s'", reportIDColumn, id)

	var report Report
	err := r.db.WithContext(ctx).Where(q).First(&report).Error
	return report.ToReportEntity(), err
}

func (r *repository) GetReportsByQuery(ctx context.Context, param entity.GetReportsByQuery) ([]*entity.Report, *pagination.Pagination, error) {
	param.ValidateQuery()

	qDB := r.db.Model(&Report{}).WithContext(ctx)

	pg := param.Pagination
	pg.ValidatePagination()

	qOrder := fmt.Sprintf("%s %s", param.Sort, param.Order)
	qQuery := fmt.Sprintf("%s LIKE %s", descriptionColumn, param.Query+"%")

	if param.Query != "" {
		qDB = qDB.Where(qQuery)
	}

	if param.Sort != "" {
		qDB = qDB.Order(qOrder)
	}

	var total int64
	qDB.Count(&total)

	pg.Total = int(total)

	var reportsModel []*Report
	if err := qDB.Offset(pg.Offset).Limit(pg.Limit).Find(&reportsModel).Error; err != nil {
		return nil, nil, err
	}

	var reports []*entity.Report
	for _, reportModel := range reportsModel {
		reports = append(reports, reportModel.ToReportEntity())
	}

	return reports, pg, nil
}

func (r *repository) Create(ctx context.Context, report *entity.Report) error {
	reportModel := Report{}.FromReportEntityToDTOWithCreatedAt(report)
	return r.db.WithContext(ctx).Create(reportModel).Error
}

func (r *repository) Update(ctx context.Context, report *entity.Report) error {
	q := fmt.Sprintf("%s = '%s'", reportIDColumn, report.ReportID)
	reportModel := Report{}.FromReportEntityToDTO(report)
	fmt.Printf("%+v", reportModel)
	return r.db.WithContext(ctx).Where(q).Updates(reportModel).Error
}

func (r *repository) Delete(ctx context.Context, id string) error {
	q := fmt.Sprintf("%s = '%s'", reportIDColumn, id)
	return r.db.WithContext(ctx).Where(q).Delete(&Report{}).Error
}
