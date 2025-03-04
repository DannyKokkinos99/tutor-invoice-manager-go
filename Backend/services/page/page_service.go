package page

import (
	"net/http"
	"tutor-invoice-manager/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewPageService(db *gorm.DB) *PageService {
	return &PageService{db: db}
}

type PageService struct {
	db *gorm.DB
}

func (s *PageService) RemoveStudent(c *gin.Context) {
	var students []models.Student
	// Fetch all students from DB
	if err := s.db.Find(&students).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error fetching students ❌")
		return
	}
	// Pass students to the remove_student.html template
	c.HTML(http.StatusOK, "remove_student.html", gin.H{"students": students})
}

func (s *PageService) EditStudent(c *gin.Context) {
	var students []models.Student
	// Fetch all students from DB
	if err := s.db.Find(&students).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error fetching students ❌")
		return
	}
	// Pass students to the edit_student.html template
	c.HTML(http.StatusOK, "edit_student.html", gin.H{"students": students})
}
