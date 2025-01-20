package db

import (
	"database/sql"
	"fmt"
	"time"
)

func TestNew(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: &UserTest{},
		Plan: &PlanTest{},
	}
}

type UserTest struct {
	ID        int
	Email     string
	FirstName string
	LastName  string
	Password  string
	Active    int
	IsAdmin   int
	CreatedAt time.Time
	UpdatedAt time.Time
	Plan      *Plan
}

func (u *UserTest) GetAll() ([]*User, error) {
	var users []*User

	user := User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password",
		Active:    1,
		IsAdmin:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	users = append(users, &user)

	return users, nil
}

func (u *UserTest) GetByEmail(email string) (*User, error) {
	user := User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password",
		Active:    1,
		IsAdmin:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return &user, nil
}

func (u *UserTest) GetOne(id int) (*User, error) {
	user := User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password",
		Active:    1,
		IsAdmin:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return &user, nil
}

func (u *UserTest) Update(user User) error {
	return nil
}

func (u *UserTest) DeleteByID(id int) error {
	return nil
}

func (u *UserTest) Insert(user User) (int, error) {
	return 1, nil
}

func (u *UserTest) ResetPassword(password string) error {
	return nil
}

func (u *UserTest) PasswordMatches(plainText string) (bool, error) {
	return true, nil
}

type PlanTest struct {
	ID                  int
	PlanName            string
	PlanAmount          int
	PlanAmountFormatted string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (p *PlanTest) GetAll() ([]*Plan, error) {
	var plans []*Plan

	plan := Plan{
		ID:                  1,
		PlanName:            "Test Plan",
		PlanAmount:          1000,
		PlanAmountFormatted: "$10.00",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	plans = append(plans, &plan)

	return plans, nil
}

func (p *PlanTest) GetOne(id int) (*Plan, error) {
	plan := Plan{
		ID:                  id,
		PlanName:            "Test Plan",
		PlanAmount:          1000,
		PlanAmountFormatted: "10.00",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	return &plan, nil
}

func (p *PlanTest) SubscribeUserToPlan(user User, plan Plan) error {
	return nil
}

func (p *PlanTest) AmountForDisplay() string {
	amount := float64(p.PlanAmount) / 100.0
	return fmt.Sprintf("$%.2f", amount)
}
