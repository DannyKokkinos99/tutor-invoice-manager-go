package routes

import (
	"tutor-invoice-manager/services/student"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterStudentRoutes(r *gin.Engine, db *gorm.DB) {
	studentService := student.NewStudentService(db)

	r.POST("/student", studentService.CreateStudent)
	r.DELETE("/student", studentService.DeleteStudent)
	r.PUT("/student", studentService.EditStudent)
	r.GET("/student", studentService.GetStudent)

}
