package model

import (
	"time"
)

// Patient represents the Patient table
type Patient struct {
	ID     *int    `json:"id"`
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	Age    *int    `json:"age"`
	Images []Image `json:"images"`
}

// Image represents the Image table
type Image struct {
	ID        *int    `json:"id"`
	PatientID *int    `json:"patient_id"`
	ImagePath *string `json:"image_path"`
	Patient   Patient `gorm:"foreignKey:PatientID;constraint:OnUpdate:CASCADE;OnDelete:SET NULL;" json:"patient"`
}

// Admin represents the Admin table
type Admin struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Username   string `gorm:"unique" json:"username"`
	Password   string `json:"password"`
	Title      string `json:"title"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
}

// Entrepreneur represents the Entrepreneur table
type Entrepreneur struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Username    string `gorm:"unique;not null;size:255" json:"username"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	Title       string `json:"title"`
	FirstName   string `json:"first_name"`
	MiddleName  string `json:"middle_name"`
	LastName    string `json:"last_name"`
	Shops       []Shop `gorm:"foreignKey:EntrepreneurID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;" json:"shops"`
}

// Shop represents the Shop table
type Shop struct {
	ID                   uint           `gorm:"primaryKey" json:"id"`
	Name                 string         `json:"name"`
	ShopCategoryID       uint           `gorm:"not null" json:"shop_category_id"`
	ShopCategory         ShopCategory   `gorm:"foreignKey:ShopCategoryID;constraint:OnDelete:SET NULL;OnUpdate:CASCADE;" json:"shop_category"`
	Status               bool           `json:"status"`
	FullDescription      string         `json:"full_description"`
	BriefDescription     string         `json:"brief_description"`
	EntrepreneurID       uint           `gorm:"not null" json:"entrepreneur_id"`
	Entrepreneur         Entrepreneur   `gorm:"foreignKey:EntrepreneurID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;" json:"entrepreneur"`
	ShopOpenDates        []ShopOpenDate `gorm:"foreignKey:ShopID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;" json:"shop_open_dates"`
	ShopMenus            []ShopMenu     `gorm:"foreignKey:ShopID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;" json:"shop_menus"`
	SocialMedia          []SocialMedia  `gorm:"foreignKey:ShopID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;" json:"social_media"`
	Photos               []Photo        `gorm:"foreignKey:ShopID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;" json:"photos"`
}

// ShopCategory represents the ShopCategory table
type ShopCategory struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Shops []Shop `gorm:"foreignKey:ShopCategoryID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;" json:"shops"`
}

// ShopOpenDate represents the ShopOpenDate table
type ShopOpenDate struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	StartTime        time.Time      `json:"start_time"`
	EndTime          time.Time      `json:"end_time"`
	ShopID           uint           `gorm:"not null" json:"shop_id"`
	MarketOpenDateID uint           `gorm:"not null" json:"market_open_date_id"`
	Shop             Shop           `gorm:"foreignKey:ShopID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;" json:"shop"`
	MarketOpenDate   MarketOpenDate `gorm:"foreignKey:MarketOpenDateID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;" json:"market_open_date"`
}

// MarketOpenDate represents the MarketOpenDate table
type MarketOpenDate struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Date          time.Time      `gorm:"type:date" json:"date"`
	StartTime     time.Time      `json:"start_time"`
	EndTime       time.Time      `json:"end_time"`
	ShopOpenDates []ShopOpenDate `gorm:"foreignKey:MarketOpenDateID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;" json:"shop_open_dates"`
}

// MarketMap represents the MarketMap table
type MarketMap struct {
	BlockID uint  `gorm:"primaryKey" json:"block_id"`
	ShopID  *uint `json:"shop_id"` 
	Shop    Shop  `gorm:"foreignKey:ShopID;constraint:OnDelete:SET NULL;OnUpdate:CASCADE;" json:"shop"`
}



// SocialMedia represents the SocialMedia table
type SocialMedia struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Platform string `json:"platform"`
	Link     string `json:"link"`
	ShopID   uint   `gorm:"not null" json:"shop_id"`
	Shop     Shop   `gorm:"foreignKey:ShopID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;" json:"shop"`
}

// ShopMenu represents the ShopMenu table
type ShopMenu struct {
	ID                 uint    `gorm:"primaryKey" json:"id"`
	ProductDescription string  `json:"product_description"`
	Price              float64 `json:"price"`
	ProductName        string  `json:"product_name"`
	ShopID             uint    `gorm:"not null" json:"shop_id"`
	Shop               Shop    `gorm:"foreignKey:ShopID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;" json:"shop"`
	Photo              Photo   `gorm:"foreignKey:MenuID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;" json:"photo"`
}

// Photo represents the Photo table
type Photo struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	PhotoCategory string `json:"photo_category"`
	PathFile      string `json:"path_file"`
	MenuID        *uint  `json:"menu_id"`
	WorkshopID    *uint  `json:"workshop_id"` // Nullable foreign key
	ShopID        *uint  `json:"shop_id"`
}

// Workshop represents the Workshop table
type Workshop struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"unique;not null" json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Language    string    `json:"language"`
	Instructor  string    `json:"instructor"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Date        time.Time `gorm:"type:date" json:"date"`
	Photos      []Photo   `gorm:"foreignKey:WorkshopID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;" json:"photos"`
}

// ContactToAdmin represents the ContactToAdmin table
type ContactToAdmin struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Problem      string `json:"problem"`
	FromUsername string `json:"from_username"`
	Detail       string `json:"detail"`
}
