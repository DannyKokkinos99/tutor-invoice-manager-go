package student

import (
	"fmt"
	"log"
	"net/http"
	"tutor-invoice-manager/models"
	"tutor-invoice-manager/schemas"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewStudentService(db *gorm.DB) *StudentService {
	return &StudentService{db: db}
}

type StudentService struct {
	db *gorm.DB
}

// CreateStudent godoc
// @Summary Create a new student
// @Description Create a new student with the provided details
// @Tags student
// @Accept json
// @Produce json
// @Param student body schemas.StudentCreate true "student details"
// @Router /student [post]
func (s *StudentService) CreateStudent(c *gin.Context) {
	var student models.Student
	if err := c.ShouldBind(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if s.db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not available ❌"})
		return
	}

	if err := s.db.Create(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("New student created - %s %s (%v)✅\n", student.Name, student.Surname, student.ID)
	c.Header("HX-Location", "/")
	c.Status(http.StatusCreated)
}

// DeleteStudent godoc
// @Summary Delete a student
// @Description Delete a student by their ID
// @Tags student
// @Param request body schemas.StudentDelete true "Student ID"
// @Router /student [delete]
func (s *StudentService) DeleteStudent(c *gin.Context) {
	var req schemas.StudentDelete
	fmt.Printf("made it into delete")
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload ❌"})
		return
	}

	var student models.Student
	if err := s.db.Where("id = ?", req.ID).First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found ❌"})
		return
	}

	if err := s.db.Delete(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Student - %s %s deleted successfully ✅\n", student.Name, student.Surname)
	c.Header("HX-Location", "/")
	c.Status(http.StatusOK)
}

// EditStudent godoc
// @Summary Update student details
// @Description Update a student's information using a JSON payload
// @Tags student
// @Param request body schemas.StudentUpdate true "Student details"
// @Router /student [put]
func (s *StudentService) EditStudent(c *gin.Context) {
	// Get raw body data for debugging
	var req schemas.StudentUpdate
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload ❌"})
		return
	}

	var student models.Student
	if err := s.db.Where("id = ?", req.ID).First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found ❌"})
		return
	}

	if err := s.db.Model(&student).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student", "details": err.Error()})
		return
	}

	log.Printf("Student - %s %s edited successfully ✅\n", student.Name, student.Surname)
	c.Header("HX-Location", "/")
	c.Status(http.StatusOK)
}

// GetStudent godoc
// @Summary Get student details
// @Description Retrieve a student's information using their ID
// @Tags student
// @Param student_id path string true "Student ID"
// @Router /student/{student_id} [get]
func (s *StudentService) GetStudent(c *gin.Context) {
	// Get student_id from the URL parameter
	studentID := c.DefaultQuery("student_id", "")
	var student models.Student
	if err := s.db.Where("id = ?", studentID).First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found ❌"})
		return
	}

	// Return student data as JSON response
	c.HTML(http.StatusOK, "filled_form.html", gin.H{"student": student})
}
