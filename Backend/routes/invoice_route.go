package routes

import (
	"tutor-invoice-manager/services/invoice"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterInvoiceRoutes(r *gin.Engine, db *gorm.DB) {
	invoiceService := invoice.NewInvoiceService(db)
	r.POST("/build_invoice", invoiceService.BuildInvoice)
	r.GET("/serve_pdf", invoiceService.ServeInvoice)
	r.GET("/invoice_send", invoiceService.SendInvoice)
	r.GET("/cancel_invoice", invoiceService.CancelInvoice)
}
