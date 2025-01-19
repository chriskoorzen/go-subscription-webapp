package main

import "github.com/chriskoorzen/go-subscription-webapp/cmd/web/db"

// For now, this is a dummy function
// In a real application, this function would generate an invoice
// using financial data from the user and plan, and accounting for taxes, discounts, etc.
func (app *Config) GenerateInvoice(u db.User, plan *db.Plan) (string, error) {
	return plan.PlanAmountFormatted, nil
}
