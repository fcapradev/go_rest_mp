package http

import (
	"net/http"
	"reliability_demo_app/internal/usecases"

	"github.com/melisource/fury_go-core/pkg/log"
	"github.com/melisource/fury_go-core/pkg/web"
)

func NewReportController(reportService usecases.ReportService) *ReportController {
	return &ReportController{reportService}
}

type ReportController struct {
	service usecases.ReportService
}

func (c *ReportController) AddRoutes(router Router) {
	router.Get("/report", c.getReport)
}

func (c *ReportController) getReport(w http.ResponseWriter, r *http.Request) error {
	result, err := c.service.GetReport(r.Context())
	if err != nil {
		log.Error(r.Context(), "get report failed", log.Err(err))
		return web.NewError(http.StatusInternalServerError, "get report failed")
	}

	return web.EncodeJSON(w, result, http.StatusOK)
}
