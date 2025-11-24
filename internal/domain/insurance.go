package domain

import (
	"time"

	"gorm.io/gorm"
)

type Insurance struct {
	gorm.Model
	Name             string  `gorm:"type:varchar(100);not null" json:"name"`
	Description      string  `gorm:"type:text" json:"description"`
	MinAge           int     `json:"min_age"`
	MaxAge           int     `json:"max_age"`
	SumAssuredLimit  float64 `gorm:"type:decimal(10,2)" json:"sum_assured_limit"`
	InsurancePremium float64 `gorm:"type:decimal(10,2)" json:"insurabce_premium"`
	Status           string  `json:"status"`
}

type InsurancePurchaseRequest struct {
	CustomerID       string `json:"customer_id" binding:"required"`
	InsuranceID      int    `json:"insurance_id" binding:"required"`
	PaymentFrequency string `json:"payment_frequency" binding:"required"`
	PaymentMethod    string `json:"payment_method" binding:"required"`
}

type CustomerInsurance struct {
	PolicyNumber     string    `gorm:"primaryKey;type:varchar(50)" json:"policy_number"`
	CustomerID       string    `gorm:"type:varchar(50);not null;index" json:"customer_id"`
	InsuranceID      int       `gorm:"not null;index" json:"insurance_id"`
	SumAssured       float64   `gorm:"type:decimal(10,2);not null" json:"sum_assured"`
	PremiumAmount    float64   `gorm:"type:decimal(10,2);not null" json:"premium_amount"`
	PaymentFrequency string    `json:"payment_frequency"`
	PaymentMethod    string    `gorm:"type:varchar(30)" json:"payment_method"`
	EffectDate       time.Time `json:"effect_date"`
	Status           string    `json:"status"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}
