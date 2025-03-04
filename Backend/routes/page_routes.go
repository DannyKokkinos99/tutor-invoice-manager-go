package routes

import (
	"net/http"
	"tutor-invoice-manager/services/page"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterPageRoutes(r *gin.Engine, db *gorm.DB) {
	pageService := page.NewPageService(db)
	r.Static("/static", "./static")
	r.LoadHTMLGlob("./templates/*")
	r.GET("/", func(c *gin.Context) {
		c.File("./templates/index.html")
	})
	r.GET("/action/add_student", func(c *gin.Context) {
		c.HTML(http.StatusOK, "add_student.html", nil)
	})
	r.GET("/action/remove_student", pageService.RemoveStudent)
	r.GET("/action/edit_student", pageService.EditStudent)
	r.GET("/action/create_invoice", pageService.CreateInvoice)
}
