package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"tutor-invoice-manager/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestFlowStudent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set up an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	for _, model := range models.AllModels {
		db.AutoMigrate(model)
	}

	r := gin.Default()
	studentService := NewStudentService(db)

	r.POST("/student", studentService.CreateStudent)
	r.DELETE("/student", studentService.DeleteStudent)
	r.PUT("/student", studentService.EditStudent)

	// TESTING CREATE STUDENT
	testStudent := `{
					"name": "John",
					"surname": "Doe",
					"parent": "Jane Doe",
					"email": "john.doe@example.com",
					"address": "123 Main St, Springfield, IL, 62701",
					"phone_number": "+1234567890",
					"price_per_hour": 50
					}`

	req, err := http.NewRequest(http.MethodPost, "/student", strings.NewReader(testStudent))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)

	// TESTING EDIT USER
	testEditStudent := `{
		"student_id": 1,
		"price_per_hour": 30
		}`
	req, err = http.NewRequest(http.MethodPut, "/student", strings.NewReader(testEditStudent))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var studentResponse struct {
		ID           uint   `json:"student_id"`
		Name         string `json:"name"`
		PricePerHour int    `json:"price_per_hour"`
	}

	err = json.NewDecoder(rr.Body).Decode(&studentResponse)
	if err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	assert.Equal(t, "John", studentResponse.Name)
	assert.Equal(t, 30, studentResponse.PricePerHour)

	// TESTING GET USER Correctly
	// req, err = http.NewRequest(http.MethodGet, "/user", nil)
	// if err != nil {
	// 	t.Fatalf("Error creating request: %v", err)
	// }
	// req.Header.Add("Authorization", token)
	// rr = httptest.NewRecorder()
	// r.ServeHTTP(rr, req)
	// assert.Equal(t, http.StatusOK, rr.Code)

}
