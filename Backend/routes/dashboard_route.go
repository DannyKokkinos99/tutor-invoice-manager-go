package routes

import (
	"tutor-invoice-manager/services/dashboard"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterDashboardRoutes(r *gin.Engine, db *gorm.DB) {
	dashboardService := dashboard.NewDashboardService(db)
	r.GET("/get_total_students", dashboardService.NumberOfStudents)
	r.GET("/get_total_hours", dashboardService.TotalHours)
	r.GET("/get_total_income", dashboardService.Income)
}
