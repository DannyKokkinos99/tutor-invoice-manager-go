package routes

import (
	"tutor-invoice-manager/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	studentService := services.NewStudentService(db)

	r.POST("/student", studentService.CreateStudent)
	r.DELETE("/student", studentService.DeleteStudent)
	r.PUT("/student", studentService.EditStudent)

}
