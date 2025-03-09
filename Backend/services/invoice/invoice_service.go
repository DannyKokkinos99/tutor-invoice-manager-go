package invoice

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"tutor-invoice-manager/gdrive"
	"tutor-invoice-manager/gmail"
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
	c.File(filePath)
}

func (s *InvoiceService) SendInvoice(c *gin.Context) {
	filePath := c.Query("invoice_name")
	studentID, _ := strconv.ParseUint(c.Query("student_id"), 10, 64)
	hours, _ := strconv.ParseFloat(c.Query("hours"), 64)
	mail, _ := strconv.ParseInt(c.Query("mail"), 10, 32)
	sendMail := mail == 1

	//get student from database
	var student models.Student
	if err := s.db.Where("id = ?", studentID).First(&student).Error; err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Student is %s\n", student.Name)
	//Save pdf to gdrive
	g, err := gdrive.NewGDrive("service_account.json")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	f, err := g.CreateFolder(student.Name, INVOICE_FOLDER_ID)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	_, err = g.UploadFile(filePath, f.Id, filePath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	//TODO: send emaail to parents including pdf
	if sendMail {

		bodyText := fmt.Sprintf(`
		<html>
			<body>
				<p>Dear %s,</p>
				<p>Please find attached the weekly invoice.</p>
				<p>Kind regards,</p>
				<p>Danny</p>
			</body>
		</html>`, student.Parent)

		fmt.Println(bodyText)
		currentDate := time.Now().Format("02-01-2006")
		subject := fmt.Sprintf("Invoice - %s", currentDate)
		gm := gmail.NewEmailSender(os.Getenv("EMAIL"), os.Getenv("APP_PASSWORD"))
		err := gm.SendEmail(student.Email, subject, bodyText, filePath)
		if err != nil {
			log.Fatalf("Error sending email: %v\n", err)
		}
	}
	// Delete local invoice after it is send via email
	err = os.Remove(filePath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	//Build invoice
	invoice := models.Invoice{
		Hours:     hours,
		Total:     hours * float64(student.PricePerHour),
		StudentID: student.ID,
		Date:      time.Now().Format("02-01-2006"),
	}
	//Update invoice counter
	student.InvoiceCount = student.InvoiceCount + 1
	s.db.Save(&student)
	//Save invoice to database
	if s.db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not available ❌"})
		return
	}
	err = s.db.Exec("SELECT setval('invoices_id_seq', (SELECT MAX(id) FROM invoices) + 1)").Error
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	if err := s.db.Create(&invoice).Error; err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("New Invoice created for %s %s (%s) ✅\n", student.Name, student.Surname, filePath)

	//TODO: change to navigate to invoice inspection step fater testing.
	c.Header("HX-Location", "/")
	c.Status(http.StatusCreated)
}

func (s *InvoiceService) BuildInvoice(c *gin.Context) {
	//Create request binding
	var req schemas.InvoiceBuilder
	if err := c.ShouldBind(&req); err != nil {
		log.Fatalf("Error: %v", err)
	}
	//get student from database
	var student models.Student
	if err := s.db.Where("id = ?", req.ID).First(&student).Error; err != nil {
		log.Fatalf("Error: %v", err)
	}

	studentFullName := fmt.Sprintf("%s %s", student.Name, student.Surname)
	invoiceName := fmt.Sprintf("Invoice-%v.pdf", student.InvoiceCount)

	//Build PDF using PDF library
	invoicePDF := pdf.NewInvoicePDF(
		req.Hours,
		float64(student.PricePerHour),
		student.InvoiceCount,
		studentFullName,
		student.Address,
		student.PhoneNumber)
	filePath := filepath.Join(INVOICE_FOLDER_LOCAL_PATH, invoiceName)
	err := invoicePDF.GeneratePDF(filePath)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)

	c.HTML(http.StatusOK, "invoice_inspection.html", gin.H{
		"student_id":   student.ID,
		"hours":        req.Hours,
		"invoice_name": invoiceName,
	})
}

func (s *InvoiceService) CancelInvoice(c *gin.Context) {
	filePath := c.Query("invoice_name")
	err := os.Remove(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete local copy of Invoice ❌"})
	}
	c.Header("HX-Location", "/")
	c.Status(http.StatusOK)
}
