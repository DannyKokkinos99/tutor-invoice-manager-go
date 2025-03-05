package invoice

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"
	"tutor-invoice-manager/models"
	"tutor-invoice-manager/pdf"
	"tutor-invoice-manager/schemas"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	INVOICE_FOLDER_ID         = "1--qhpO7fr5q4q7x0pRxdiETcFyBsNOGN"
	INVOICE_FOLDER_LOCAL_PATH = "./"
)

func NewInvoiceService(db *gorm.DB) *InvoiceService {
	return &InvoiceService{db: db}
}

type InvoiceService struct {
	db *gorm.DB
}

func (s *InvoiceService) ServeInvoice(c *gin.Context) {
	filePath := c.Query("invoice_name")
	if filePath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Missing invoice_name in query params ❌"})
		return
	}
	fmt.Printf("File path : %s", filePath)
	c.File(filePath)
}

func (s *InvoiceService) BuildInvoice(c *gin.Context) {
	date := time.Now().Format("02-01-2006")
	//Create request binding
	var req schemas.InvoiceBuilder
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//get student from database
	var student models.Student
	if err := s.db.Where("id = ?", req.ID).First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found ❌"})
		return
	}
	fmt.Printf("Student found: %v (%v)\n ", student, student.ID)
	//Build invoice
	invoice := models.Invoice{
		Hours:     req.Hours,
		Total:     req.Hours * float64(student.PricePerHour),
		StudentID: student.ID,
		Date:      date,
	}
	fmt.Printf("Invoice created: %v\n", invoice)

	//Determine the invoiceNumber by querying the database
	var tempInvoices []models.Invoice
	if err := s.db.Where("student_id = ?", req.ID).Find(&tempInvoices).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No invoices found for this student ❌"})
		return
	}
	fmt.Printf("Number of invoices the student has %v\n", len(tempInvoices))
	nextInvoiceIndex := len(tempInvoices) + 1
	studentFullName := fmt.Sprintf("%s %s", student.Name, student.Surname)
	invoiceName := fmt.Sprintf("Invoice-%v.pdf", nextInvoiceIndex)

	//Build PDF using PDF library
	invoicePDF := pdf.NewInvoicePDF(
		req.Hours,
		float64(student.PricePerHour),
		nextInvoiceIndex,
		studentFullName,
		student.Address,
		student.PhoneNumber)

	filePath := filepath.Join(INVOICE_FOLDER_LOCAL_PATH, invoiceName)
	fmt.Printf("INVOICE PATH %s\n", filePath)

	err := invoicePDF.GeneratePDF(filePath)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	// GO TO INSPECTION PAGE
	c.HTML(http.StatusOK, "invoice_inspection.html", gin.H{
		"student_id":   student.ID,
		"hours":        req.Hours,
		"total":        req.Hours * float64(student.PricePerHour),
		"invoice_name": invoiceName,
	})

	// //Save pdf to gdrive
	// g, err := gdrive.NewGDrive("service_account.json")
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Google Drive ❌"})
	// }
	// f, err := g.CreateFolder(student.Name, INVOICE_FOLDER_ID)
	// if err != nil {
	// 	fmt.Println(err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Google Drive ❌"})
	// }
	// _, err = g.UploadFile(filePath, f.Id, invoiceName)
	// if err != nil {
	// 	fmt.Println(err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to Google Drive ❌"})
	// }

	//TODO: send emaail to parents including pdf

	// // Delete local invoice after it is send via email
	// err = os.Remove(filePath)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete local copy of Invoice ❌"})
	// }

	// //Save invoice to database
	// if s.db == nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not available ❌"})
	// 	return
	// }

	// if err := s.db.Create(&invoice).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	// fmt.Printf("New Invoice created - %s %s (Invoice-%v.pdf) ✅\n", student.Name, student.Surname, nextInvoiceIndex)

	//TODO: change to navigate to invoice inspection step fater testing.
	c.Header("HX-Location", "/")
	c.Status(http.StatusCreated)
}
