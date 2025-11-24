package domain

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	CustomerID string `gorm:"primaryKey" json:"customer_id"`
	Firstname  string `json:"firstname" binding:"required"`
	Lastname   string `json:"lastname" binding:"required"`
	Birthdate  string `json:"birthdate"`
	Age        int    `json:"age" binding:"required,gt=0"`
	Gender     string `json:"gender" binding:"required,oneof=Male Female Other"`
	Email      string `json:"email" gorm:"unique"`
	Phone      string `json:"phone" gorm:"unique"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type Beneficiaries struct {
	gorm.Model
	CustomerID   string  `gorm:"type:varchar(50);not null;index" json:"customer_id"`
	InsuranceID  int     `gorm:"index" json:"insurance_id"`
	PolicyNumber string  `gorm:"type:varchar(20);index" json:"policy_number"`
	Firstname    string  `json:"firstname"`
	Lastname     string  `json:"lastname"`
	Relationship string  `json:"relationship"`
	Percentage   float64 `gorm:"type:decimal(5,2)" json:"percentage"`
	Email        string  `json:"email"`
	Phone        string  `json:"phone"`
}
