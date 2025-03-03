package services

import (
	"fmt"
	"log"
	"net/http"
	"tutor-invoice-manager/models"
	"tutor-invoice-manager/schemas"
	"tutor-invoice-manager/utils"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
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
	if err := c.ShouldBindJSON(&student); err != nil {
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
	studentGet := schemas.StudentGet{}
	err := mapstructure.Decode(student, &studentGet)
	if err != nil {
		fmt.Printf("Error decoding ❌: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to map response ❌"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": studentGet})
}

// DeleteStudent godoc
// @Summary Delete a student
// @Description Delete a student by their ID
// @Tags student
// @Param request body schemas.StudentDelete true "Student ID"
// @Router /student [delete]
func (s *StudentService) DeleteStudent(c *gin.Context) {
	var req schemas.StudentDelete

	if err := c.ShouldBindJSON(&req); err != nil {
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
	message := fmt.Sprintf("Student - %s %s deleted successfully ✅", student.Name, student.Surname)
	log.Println(message)
	c.JSON(http.StatusOK, gin.H{"message": message})
}

// EditStudent godoc
// @Summary Update student details
// @Description Update a student's information using a JSON payload
// @Tags student
// @Param request body schemas.StudentUpdate true "Student details"
// @Router /student [put]
func (s *StudentService) EditStudent(c *gin.Context) {
	var req schemas.StudentUpdate

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload ❌"})
		return
	}

	var student models.Student
	if err := s.db.Where("id = ?", req.ID).First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found ❌"})
		return
	}

	updateData := utils.EditObject(req)
	if err := s.db.Model(&student).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	studentGet := schemas.StudentGet{}
	err := mapstructure.Decode(student, &studentGet)
	if err != nil {
		fmt.Printf("Error decoding ❌: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to map response ❌"})
		return
	}
	c.JSON(http.StatusOK, studentGet)
}
