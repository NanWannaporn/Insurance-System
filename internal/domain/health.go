package domain

import "gorm.io/gorm"

type HealthDeclaration struct {
	gorm.Model
	CustomerID string  `gorm:"type:varchar(50);not null;uniqueIndex" json:"customer_id"`
	Height     float64 `gorm:"type:decimal(5,2)" json:"height" binding:"required"`
	Weight     float64 `gorm:"type:decimal(5,2)" json:"weight" binding:"required"`
	BloodGroup string  `json:"blood_group" binding:"required"`
	//Medical History)
	HasChronicDisease           bool   `json:"has_chronic_disease"`
	ChronicDisease              string `json:"chronic_disease"`
	HasBeenHospitalizedLastYear bool   `json:"has_been_hospitalized_last_year"`
	SmokingStatus               string `json:"smoking_status"`
	MedicalDetails              string `json:"medical_details"`
	SurgicalHistory             string `json:"surgical_history"`
	Allergies                   string `json:"allergies"`
	FamilyMedicalHistory        string `json:"family_medical_history"`
}
