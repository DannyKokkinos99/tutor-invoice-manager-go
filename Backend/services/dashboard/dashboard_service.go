package dashboard

import (
	"net/http"
	"tutor-invoice-manager/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{db: db}
}

type DashboardService struct {
	db *gorm.DB
}

func (d *DashboardService) NumberOfStudents(c *gin.Context) {
	// Query the total number of students in the database
	var count int64
	if err := d.db.Model(&models.Student{}).Count(&count).Error; err != nil {
		// If there is an error, return a 500 status with the error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get the number of students ❌"})
		return
	}
	c.String(http.StatusOK, "%d", count)
}

func (d *DashboardService) TotalHours(c *gin.Context) {
	// Query the total number of hours in the invoices table
	var totalHours float64
	if err := d.db.Model(&models.Invoice{}).Select("SUM(hours)").Scan(&totalHours).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		totalHours = 0
	}

	// Return the total hours as a string (or you could return it as JSON, depending on your need)
	c.String(http.StatusOK, "%.2f", totalHours)
}

func (d *DashboardService) Income(c *gin.Context) {
	// Query the total number of hours in the invoices table
	var totalIncome float64
	if err := d.db.Model(&models.Invoice{}).Select("SUM(total)").Scan(&totalIncome).Error; err != nil {
		totalIncome = 0
	}

	// Return the total hours as a string (or you could return it as JSON, depending on your need)
	c.String(http.StatusOK, "£%.0f", totalIncome)
}
