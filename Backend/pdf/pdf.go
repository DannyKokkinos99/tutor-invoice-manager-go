package pdf

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

const (
	POUND_SYMBOL = "£"
)

// InvoicePDF represents an invoice PDF
type InvoicePDF struct {
	pdf              *gofpdf.Fpdf
	Quantity         float64
	UnitPrice        float64
	InvoiceNumber    int
	CustomerFullName string
	CustomerAddress  string
	CustomerPhone    string
}

// NewInvoicePDF initializes a new invoice PDF
func NewInvoicePDF(
	quantity float64,
	unitPrice float64,
	invoiceNumber int,
	customerFullName string,
	customerAddress string,
	customerPhone string,
) *InvoicePDF {
	pdf := gofpdf.New("P", "mm", "A4", "")

	return &InvoicePDF{
		pdf:              pdf,
		Quantity:         quantity,
		UnitPrice:        unitPrice,
		InvoiceNumber:    invoiceNumber,
		CustomerFullName: customerFullName,
		CustomerAddress:  customerAddress,
		CustomerPhone:    customerPhone,
	}
}

// Header creates the invoice header
func (i *InvoicePDF) Header() {
	i.pdf.SetFont("Arial", "B", 24)
	i.pdf.CellFormat(0, 10, "INVOICE", "", 1, "C", false, 0, "")
	i.pdf.Ln(10)
}

// CompanyDetails adds company information
func (i *InvoicePDF) CompanyDetails() {
	x, y := 10.0, 30.0
	i.pdf.SetFont("Arial", "B", 12)
	i.pdf.SetXY(x, y)
	i.pdf.CellFormat(90, 35, "", "1", 0, "", false, 0, "")

	i.pdf.SetFillColor(175, 238, 238)
	i.pdf.SetXY(x, y)
	i.pdf.CellFormat(90, 10, "", "1", 0, "", true, 0, "")

	i.pdf.SetXY(x, y)
	i.pdf.Cell(75, 10, "From")
	i.pdf.SetXY(x, y+10)
	i.pdf.Cell(75, 10, "Ntani Kokkinos")
	i.pdf.SetXY(x, y+18)
	i.pdf.MultiCell(80, 5, "St. Thomas Street, The Milliners 308, BS1 6WT", "", "L", false)
	i.pdf.SetXY(x, y+26)
	i.pdf.Cell(75, 10, "+44 7448 646758")
	i.pdf.Ln(20)
}

// CustomerDetails adds customer information
func (i *InvoicePDF) CustomerDetails(name, address, phoneNumber string) {
	x, y := 10.0, 80.0

	i.pdf.SetFont("Arial", "B", 12)
	i.pdf.SetXY(x, y)
	i.pdf.CellFormat(90, 35, "", "1", 0, "", false, 0, "")

	i.pdf.SetFillColor(175, 238, 238)
	i.pdf.SetXY(x, y)
	i.pdf.CellFormat(90, 10, "", "1", 0, "", true, 0, "")

	i.pdf.SetXY(x, y)
	i.pdf.Cell(75, 10, "Bill to")
	i.pdf.SetXY(x, y+10)
	i.pdf.Cell(75, 10, fmt.Sprint(name))
	i.pdf.SetXY(x, y+18)
	i.pdf.MultiCell(80, 5, address, "", "L", false)
	i.pdf.SetXY(x, y+26)
	i.pdf.Cell(75, 10, phoneNumber)
	i.pdf.Ln(20)
}

// InvoiceDetails adds invoice metadata
func (i *InvoicePDF) InvoiceDetails(invoiceNumber int) {
	today := time.Now().Format("02/01/2006")

	x, y := 100.0, 30.0
	i.pdf.SetFont("Arial", "", 10)

	i.pdf.SetXY(x, y)
	i.pdf.CellFormat(80, 10, fmt.Sprintf("Invoice Date: %s", today), "", 0, "R", false, 0, "")

	i.pdf.SetXY(x, y+10)
	i.pdf.CellFormat(80, 10, "Due: On Receipt", "", 0, "R", false, 0, "")

	i.pdf.SetXY(x, y+20)
	i.pdf.CellFormat(80, 10, fmt.Sprintf("Invoice Number: %d", invoiceNumber), "", 0, "R", false, 0, "")
}

// InvoiceTable creates the table with quantity, unit price, and total
func (i *InvoicePDF) InvoiceTable(quantity float64, unitPrice, total float64) {
	tr := i.pdf.UnicodeTranslatorFromDescriptor("")
	x, y := 10.0, 135.0
	i.pdf.SetFont("Arial", "B", 12)
	i.pdf.SetFillColor(175, 238, 238)

	i.pdf.SetXY(x, y)
	i.pdf.CellFormat(80, 10, "Item", "1", 0, "C", true, 0, "")
	i.pdf.CellFormat(30, 10, "Quantity", "1", 0, "C", true, 0, "")
	i.pdf.CellFormat(40, 10, "Unit Price", "1", 0, "C", true, 0, "")
	i.pdf.CellFormat(40, 10, "Total", "1", 0, "C", true, 0, "")
	i.pdf.Ln(10)

	i.pdf.SetFont("Arial", "", 12)
	i.pdf.CellFormat(80, 10, "Tuition Service", "1", 0, "C", false, 0, "")
	i.pdf.CellFormat(30, 10, fmt.Sprintf("%v", quantity), "1", 0, "C", false, 0, "")
	i.pdf.CellFormat(40, 10, fmt.Sprintf("%s%.2f", tr(POUND_SYMBOL), unitPrice), "1", 0, "C", false, 0, "")
	i.pdf.CellFormat(40, 10, fmt.Sprintf("%s%.2f", tr(POUND_SYMBOL), total), "1", 0, "C", false, 0, "")
	i.pdf.Ln(10)
}

// TotalAmount adds the total amount at the bottom
func (i *InvoicePDF) TotalAmount(total float64) {
	tr := i.pdf.UnicodeTranslatorFromDescriptor("")
	i.pdf.SetFont("Arial", "B", 12)
	i.pdf.CellFormat(0, 10, fmt.Sprintf("Total: %s%.2f", tr(POUND_SYMBOL), total), "", 1, "R", false, 0, "")
	i.pdf.SetXY(10, 180)
	i.pdf.Cell(0, 10, "INVOICE IS PAID")
}

// GeneratePDF creates the full invoice and saves it
func (i *InvoicePDF) GeneratePDF(filename string) error {
	i.pdf.AddPage()
	i.Header()
	i.CompanyDetails()
	i.CustomerDetails(i.CustomerFullName, i.CustomerAddress, i.CustomerPhone)
	i.InvoiceDetails(i.InvoiceNumber)
	total := i.Quantity * i.UnitPrice
	i.InvoiceTable(i.Quantity, i.UnitPrice, total)
	i.TotalAmount(total)

	return i.pdf.OutputFileAndClose(filename)
}
