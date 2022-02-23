package http

import (
	"context"
	"net/http"

	"github.com/Aldiwildan77/backend-hexa-template/core/entity"
	"github.com/Aldiwildan77/backend-hexa-template/core/module"
	"github.com/Aldiwildan77/backend-hexa-template/pkg/convert"
	"github.com/Aldiwildan77/backend-hexa-template/pkg/pagination"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ReportHandler interface {
	CreateReport(e echo.Context) error
	UpdateReport(e echo.Context) error
	DeleteReport(e echo.Context) error
	GetReport(e echo.Context) error
	GetReports(e echo.Context) error
}

type reportHandler struct {
	uc module.ReportUsecase
}

func NewReportHandler(uc module.ReportUsecase) ReportHandler {
	return &reportHandler{uc: uc}
}

func (h *reportHandler) CreateReport(c echo.Context) error {
	ctx := context.Background()
	vld := validator.New()

	var report *entity.Report
	if err := c.Bind(&report); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	}

	if err := vld.Struct(report); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	}

	report, err := h.uc.CreateReport(ctx, report)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, Response{Data: report, Message: "create success"})
}

func (h *reportHandler) UpdateReport(e echo.Context) error {
	ctx := context.Background()
	vld := validator.New()

	var report *entity.Report
	if err := e.Bind(&report); err != nil {
		return e.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	}

	report.ReportID = e.Param("id")

	if err := vld.Struct(report); err != nil {
		return e.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	}

	report, err := h.uc.UpdateReport(ctx, report)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err)
	}

	return e.JSON(http.StatusOK, Response{Data: report, Message: "update success"})
}

func (h *reportHandler) DeleteReport(e echo.Context) error {
	ctx := context.Background()

	if err := h.uc.DeleteReport(ctx, e.Param("id")); err != nil {
		return e.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}

	return e.JSON(http.StatusNoContent, nil)
}

func (h *reportHandler) GetReport(e echo.Context) error {
	ctx := context.Background()

	report, err := h.uc.GetReport(ctx, e.Param("id"))
	if err != nil {
		return e.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, Response{Data: report, Message: "fetch success"})
}

func (h *reportHandler) GetReports(e echo.Context) error {
	ctx := context.Background()

	param := entity.GetReportsByQuery{
		Pagination: &pagination.Pagination{
			Limit: convert.StringToInt(e.QueryParam("limit"), 10),
			Page:  convert.StringToInt(e.QueryParam("page"), 1),
		},
		Order: e.QueryParam("order"),
		Sort:  e.QueryParam("sort"),
		Query: e.QueryParam("query"),
	}

	reports, pagination, err := h.uc.GetReports(ctx, param)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, Response{
		Data:    reports,
		Message: "fetch success",
		MetadataResponse: &MetadataResponse{
			Pagination: pagination,
		},
	})
}
