package main

import (
	"time"

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

// Admin represents the Admin table
type Admin struct {
	gorm.Model
	Username   string `gorm:"primaryKey"`
	Password   string
	Title      string
	FirstName  string
	MiddleName string
	LastName   string
}

// Entrepreneur represents the Entrepreneur table
type Entrepreneur struct {
	gorm.Model
	Username   string `gorm:"primaryKey"`
	Password   string
	PhoneNumber string
	Title      string
	FirstName  string
	MiddleName string
	LastName   string
	Shops      []Shop `gorm:"foreignKey:EntUsername"`
}

// Shop represents the Shop table
type Shop struct {
	gorm.Model
	ID              uint           `gorm:"primaryKey"`
	Name            string
	ShopCategoryID  uint           `gorm:"not null"`
	ShopCategory    ShopCategory   `gorm:"foreignKey:ShopCategoryID"`
	Status          bool
	FullDescription string
	BriefDescription string
	EntUsername     string         `gorm:"not null"`
	Entrepreneur    Entrepreneur   `gorm:"foreignKey:EntUsername"`
	ShopOpenDates   []ShopOpenDate `gorm:"foreignKey:ShopID"`
	ShopMenus       []ShopMenu     `gorm:"foreignKey:ShopID"`
	SocialMedia     []SocialMedia  `gorm:"foreignKey:ShopID"`
	Photos          []Photo        `gorm:"foreignKey:ShopID"`
}

// ShopCategory represents the Shop_Category table
type ShopCategory struct {
	gorm.Model
	ID   uint   `gorm:"primaryKey"`
	Name string
	Shops []Shop `gorm:"foreignKey:ShopCategoryID"`
}

// ShopOpenDate represents the ShopOpenDate table
type ShopOpenDate struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey"`
	StartTime time.Time
	EndTime   time.Time
	ShopID    uint      `gorm:"not null"`
	DateID    uint      `gorm:"not null"`
	Shop      Shop      `gorm:"foreignKey:ShopID"`
}

// MarketOpenDate represents the MarketOpenDate table
type MarketOpenDate struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey"`
	Date      time.Time `gorm:"type:date"`
	StartTime time.Time
	EndTime   time.Time
}

// MarketMap represents the MarketMap table
type MarketMap struct {
	gorm.Model
	BlockID uint `gorm:"primaryKey"`
	ShopID  uint `gorm:"not null"`
	Shop    Shop `gorm:"foreignKey:ShopID"`
}

// SocialMedia represents the SocialMedia table
type SocialMedia struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Platform string
	Link     string
	ShopID   uint   `gorm:"not null"`
	Shop     Shop   `gorm:"foreignKey:ShopID"`
}

// ShopMenu represents the ShopMenu table
type ShopMenu struct {
	gorm.Model
	ID               uint    `gorm:"primaryKey"`
	ProductDescription string
	Price            float64
	ProductName      string
	ShopID           uint    `gorm:"not null"`
	Shop             Shop    `gorm:"foreignKey:ShopID"`
	Photos           []Photo `gorm:"foreignKey:ProductID"`
}

// Photo represents the Photo table
type Photo struct {
	gorm.Model
	ID           uint   `gorm:"primaryKey"`
	PhotoCategory string
	PathFile      string
	ProductID     *uint
	WorkshopName  *string
	ShopID        *uint
}

// Workshop represents the Workshop table
type Workshop struct {
	gorm.Model
	Name        string    `gorm:"primaryKey"`
	Description string
	Price       float64
	Language    string
	Instructor  string
	StartTime   time.Time
	EndTime     time.Time
	Date 		time.Time `gorm:"type:date"`
	Photos      []Photo `gorm:"foreignKey:WorkshopName"`
}

// ContactToAdmin represents the ContactToAdmin table
type ContactToAdmin struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	Problem     string
	FromUsername string
	Detail      string
}

