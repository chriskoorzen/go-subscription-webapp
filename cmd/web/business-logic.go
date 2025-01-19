package main

import (
	"fmt"

	"github.com/chriskoorzen/go-subscription-webapp/cmd/web/db"
	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"
)

// For now, this is a dummy function
// In a real application, this function would generate an invoice
// using financial data from the user and plan, and accounting for taxes, discounts, etc.
func (app *Config) GenerateInvoice(u db.User, plan *db.Plan) (string, error) {
	return plan.PlanAmountFormatted, nil
}

// For now, this is a dummy function
// In a real application, this function would generate a custom manual pdf document
func (app *Config) GenerateManual(u db.User, plan *db.Plan) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 15, 10)

	pdfImporter := gofpdi.NewImporter()

	pdfTemplate := pdfImporter.ImportPage(pdf, "./pdf/manual.pdf", 1, "/MediaBox")

	pdf.AddPage()

	pdfImporter.UseImportedTemplate(pdf, pdfTemplate, 0, 0, 210, 0) // 210mm x 297mm

	pdf.SetX(75)  // 75mm from left
	pdf.SetY(150) // 150mm from top

	pdf.SetFont("Arial", "", 12)                                                       // Arial, regular, 12pt
	pdf.MultiCell(0, 4, fmt.Sprintf("%s %s", u.FirstName, u.LastName), "", "C", false) // 0 width, 4 height, center align, no fill
	pdf.Ln(6)                                                                          // line break
	pdf.MultiCell(0, 4, fmt.Sprintf("%s User Guide", plan.PlanName), "", "C", false)

	return pdf
}
