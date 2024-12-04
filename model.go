package main

import (
	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	ID *int `json:"ID"`
	Name *string `json:"Name"`
	Email *string `json:"Email"`
	Age *int `json:"Age"`
	Images []Image 
}

type Image struct {
	gorm.Model
	ID        *int    `json:"ID"`
	PatientID *int    `json:"PatientID"`
	ImagePath *string `json:"ImagePath"`
	Patient   Patient `gorm:"foreignKey:PatientID"` // Establish relationship with Patient
}


