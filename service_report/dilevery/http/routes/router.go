package report_route

import (
	"github.com/gin-gonic/gin"
	report_handler "github.com/ryvasa/go-super-farmer/service_report/dilevery/http/handler"
)

type Router interface {
}

func NewRoutes(handlers *report_handler.Handlers) *gin.Engine {
	r := gin.Default()

	NewReportRoutes(r, handlers.ReportHandler)

	return r
}
