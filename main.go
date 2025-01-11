package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)



func main() {
	// db connection section
	db := SetupDatabase()

	// for development only
	// db.AutoMigrate(&Patient{}, &Image{})
	 db.AutoMigrate(&Admin{},
		&Patient{}, 
		&Image{},
		&Entrepreneur{},
		&Shop{},
		&ShopCategory{},
		&ShopOpenDate{},
		&MarketOpenDate{},
		&MarketMap{},
		&SocialMedia{},
		&ShopMenu{},
		&Photo{},
		&Workshop{},
		&ContactToAdmin{},)

	
	// use godotenv to get .env variables
	if err := godotenv.Load(); err!= nil { // gogotenv init
		log.Fatal("load .env error")
	}

	createUploadsDirectory()

	// ตัวแทนการสื่อสารกับ http server
	app := fiber.New() // fiber init

	// Static file serving to access uploaded images via /upload/{filename}
	app.Static("/upload", "./uploads")

	// Apply CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, http://127.0.0.1:3000, http://172.18.0.4:3000, http://localhost ",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowCredentials: true, // If you are using cookies or credentials
	}))

	app.Options("*", func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		return c.SendStatus(fiber.StatusNoContent)
	})
	

						// c = response and request fiber context
	app.Get("/hello" , func(c *fiber.Ctx) error {return c.SendString("test gogo")})

	// test route 
	app.Get("/patient", func(c *fiber.Ctx) error {
		return getPatients(db, c)
	  })
	app.Get("/patient/:id", func(c *fiber.Ctx) error {
		return getPatientID(db, c)
	  })
	app.Post("/patient", func(c *fiber.Ctx) error {
		return createPatient(db, c)
	  })
	app.Put("/patient/:id", func(c *fiber.Ctx) error {
		return updatePatient(db, c)
	  })
	app.Delete("/patient/:id", func(c *fiber.Ctx) error {
		return deletePatient(db, c)
	  })
	  app.Post("/patient/:id/images", func(c *fiber.Ctx) error {
		return uploadImage(db, c)
	})
	app.Get("/patient/:id/images", func(c *fiber.Ctx) error {
		return getPatientImages(db, c)
	})
	
	
	app.Get("/config", getENV)

	app.Listen(":8080")

}

func ptrString(s string) *string {
	return &s
}

// func ptrInt(i int) *int {
// 	return &i
// }

func getENV (c *fiber.Ctx) error {
	
	secret := os.Getenv("SECRET")
	if secret == ""{
		secret = "defaultsecret"
	}
	return c.JSON(fiber.Map{ 
		"SECRET": secret,
	})

}

func SetupDatabase() *gorm.DB {
	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
		  SlowThreshold: time.Second, // Slow SQL threshold
		  LogLevel:      logger.Info, // Log level
		  Colorful:      true,        // Enable color
		},
	)

	// Get DSN from environment or use default for development
	dsn := os.Getenv("DB_DSN")
	fmt.Print(dsn)
	if dsn == "" {
		// Default for development
		dsn = "user:12345678@tcp(127.0.0.1:3306)/medic?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger,})
	  
	if err != nil {
	  log.Fatal(err)
	}
  
	
	fmt.Println("Successfully DB connected!")
  
	return db
  }

  func createUploadsDirectory() {
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		os.Mkdir("./uploads", os.ModePerm)
	}
}