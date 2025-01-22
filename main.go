package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/HealthMe-pls/medic-go-api/controller"
	"github.com/HealthMe-pls/medic-go-api/model"
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
	if err := db.AutoMigrate(
		&model.Patient{},
		&model.Image{},
		&model.Admin{},
		&model.ContactToAdmin{},
		&model.ShopCategory{},
		&model.MarketOpenDate{},
		&model.Entrepreneur{},
		&model.Shop{},
		&model.ShopOpenDate{},
		&model.MarketMap{},
		&model.SocialMedia{},
		&model.ShopMenu{},
		&model.Workshop{},
		&model.Photo{}); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	// use godotenv to get .env variables
	if err := godotenv.Load(); err != nil { // gogotenv init
		log.Fatal("load .env error")
	}

	createUploadsDirectory()

	// ตัวแทนการสื่อสารกับ http server
	app := fiber.New(fiber.Config{
		JSONDecoder: json.Unmarshal,
	})
	// fiber init

	// Static file serving to access uploaded images via /upload/{filename}
	app.Static("/upload", "./uploads")

	// Apply CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://127.0.0.1:3000, http://172.18.0.4:3000, http://localhost ",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true, // If you are using cookies or credentials
	}))

	app.Options("*", func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		return c.SendStatus(fiber.StatusNoContent)
	})

	// c = response and request fiber context
	app.Get("/hello", func(c *fiber.Ctx) error { return c.SendString("test gogo") })

	// test route
	app.Get("/patient", func(c *fiber.Ctx) error { return controller.GetPatients(db, c) })
	app.Get("/patient/:id", func(c *fiber.Ctx) error { return controller.GetPatientID(db, c) })
	app.Post("/patient", func(c *fiber.Ctx) error { return controller.CreatePatient(db, c) })
	app.Put("/patient/:id", func(c *fiber.Ctx) error { return controller.UpdatePatient(db, c) })
	app.Delete("/patient/:id", func(c *fiber.Ctx) error { return controller.DeletePatient(db, c) })
	app.Post("/patient/:id/images", func(c *fiber.Ctx) error { return controller.UploadImage(db, c) })
	app.Get("/patient/:id/images", func(c *fiber.Ctx) error { return controller.GetPatientImages(db, c) })

	//admin --check
	app.Get("/admin", func(c *fiber.Ctx) error { return controller.GetAdmins(db, c) })
	app.Get("/admin/:id", func(c *fiber.Ctx) error { return controller.GetAdminByUsername(db, c) })
	app.Post("/admin", func(c *fiber.Ctx) error { return controller.CreateAdmin(db, c) })
	app.Put("/admin/:id", func(c *fiber.Ctx) error { return controller.UpdateAdmin(db, c) })
	app.Delete("/admin/:id", func(c *fiber.Ctx) error { return controller.DeleteAdmin(db, c) })

	//entrepreneur
	app.Get("/entrepreneur", func(c *fiber.Ctx) error { return controller.GetEntrepreneur(db, c) })
	app.Get("/entrepreneur/:id", func(c *fiber.Ctx) error { return controller.GetEntrepreneurByID(db, c) })
	app.Post("/entrepreneur", func(c *fiber.Ctx) error { return controller.CreateEntrepreneur(db, c) })
	app.Put("/entrepreneur/:id", func(c *fiber.Ctx) error { return controller.UpdateEntrepreneur(db, c) })
	//Not available
	app.Delete("/entrepreneur/:id", func(c *fiber.Ctx) error { return controller.DeleteEntrepreneurAndShops(db, c) })

	//map --check
	app.Get("/map", func(c *fiber.Ctx) error { return controller.GetMarketMap(db, c) })
	app.Get("/mapdetail", func(c *fiber.Ctx) error { return controller.GetMarketMapDetail(db, c) })
	app.Get("/map/:id", func(c *fiber.Ctx) error { return controller.GetMapByBlockID(db, c) })
	app.Get("/shopInmap/:id", func(c *fiber.Ctx) error { return controller.GetShopInMapID(db, c) })
	app.Post("/map", func(c *fiber.Ctx) error { return controller.CreateMarketMap(db, c) })
	app.Delete("/map/:block_id", func(c *fiber.Ctx) error { return controller.DeleteMarketMapsByBlockID(db, c) })
	app.Put("/map/:block_id", func(c *fiber.Ctx) error { return controller.UpdateMarketMapByBlockID(db, c) })
	//name
	app.Get("/mapN/:block_name", func(c *fiber.Ctx) error { return controller.GetMapByBlockName(db, c) })
	app.Delete("/mapN/:block_name", func(c *fiber.Ctx) error { return controller.DeleteMarketMapsByBlockName(db, c) })
	app.Put("/mapN/:block_name", func(c *fiber.Ctx) error { return controller.UpdateMarketMapByBlockName(db, c) })

	//shop category
	app.Post("/shopcategory", func(c *fiber.Ctx) error { return controller.CreateShopCategory(db, c) })
	app.Get("/shopcategory", func(c *fiber.Ctx) error { return controller.GetShopCategories(db, c) })
	app.Get("/shopcategory/:id", func(c *fiber.Ctx) error { return controller.GetShopCategoryByID(db, c) })
	app.Delete("/shopcategory/:id", func(c *fiber.Ctx) error { return controller.DeleteShopCategory(db, c) })
	app.Put("/shopcategory/:id", func(c *fiber.Ctx) error { return controller.UpdateShopCategory(db, c) })

	//shop
	app.Post("/shop", func(c *fiber.Ctx) error { return controller.CreateShop(db, c) })
	app.Get("/shop/:id", func(c *fiber.Ctx) error { return controller.GetShopByID(db, c) })
	// app.Get("/shopdetail/:id", func(c *fiber.Ctx) error {return controller.GetShopDetail(db, c)})
	app.Get("/shop", func(c *fiber.Ctx) error { return controller.GetShops(db, c) })
	app.Put("/shop/:id", func(c *fiber.Ctx) error { return controller.UpdateShop(db, c) })
	app.Delete("/shop/:id", func(c *fiber.Ctx) error { return controller.DeleteShop(db, c) })
	app.Get("/shops/category/:shop_category_id", func(c *fiber.Ctx) error { return controller.GetShopsByCategory(db, c) })

	// Workshop Routes
	app.Get("/workshops", func(c *fiber.Ctx) error { return controller.GetWorkshops(db, c) })
	app.Get("/workshops/:id", func(c *fiber.Ctx) error { return controller.GetWorkshopByID(db, c) })
	app.Post("/workshops", func(c *fiber.Ctx) error { return controller.CreateWorkshop(db, c) })
	app.Put("/workshops/:id", func(c *fiber.Ctx) error { return controller.UpdateWorkshop(db, c) })
	app.Delete("/workshops/:id", func(c *fiber.Ctx) error { return controller.DeleteWorkshop(db, c) })

	//manage
	app.Post("/marketDate", func(c *fiber.Ctx) error { return controller.CreateMarketOpenDate(db, c) })
	app.Get("/marketDate/:id", func(c *fiber.Ctx) error { return controller.GetMarketOpenDate(db, c) })
	app.Put("/marketDate/:id", func(c *fiber.Ctx) error { return controller.UpdateMarketOpenDate(db, c) })
	app.Delete("/marketDate/:id", func(c *fiber.Ctx) error { return controller.DeleteMarketOpenDate(db, c) })

	//social media
	app.Post("/social", func(c *fiber.Ctx) error { return controller.CreateSocialMedia(db, c) })
	app.Get("/social/:id", func(c *fiber.Ctx) error { return controller.GetSocialMedia(db, c) })
	app.Get("/social/shop/:shop_id", func(c *fiber.Ctx) error { return controller.GetSocialMediaByShopID(db, c) })
	app.Put("/social/:id", func(c *fiber.Ctx) error { return controller.UpdateSocialMedia(db, c) })
	app.Delete("/social/:id", func(c *fiber.Ctx) error { return controller.DeleteSocialMedia(db, c) })

	//shop open time
	app.Post("/shoptime", func(c *fiber.Ctx) error { return controller.CreateShopOpenDate(db, c) })
	app.Get("/shoptime/:id", func(c *fiber.Ctx) error { return controller.GetShopOpenDate(db, c) })
	app.Get("/shoptime/shop/:shop_id", func(c *fiber.Ctx) error { return controller.GetShopOpenDateByShopID(db, c) })
	app.Put("/shoptime/:id", func(c *fiber.Ctx) error { return controller.UpdateShopOpenDate(db, c) })
	app.Delete("/shoptime/:id", func(c *fiber.Ctx) error { return controller.DeleteShopOpenDate(db, c) })

	//shop menu
	app.Post("/shopmenu", func(c *fiber.Ctx) error { return controller.CreateShopMenu(db, c) })
	app.Get("/shopmenu/:id", func(c *fiber.Ctx) error { return controller.GetShopMenu(db, c) })
	app.Get("/shopmenu/shop/:shop_id", func(c *fiber.Ctx) error { return controller.GetShopMenuByShopID(db, c) })
	app.Put("/shopmenu/:id", func(c *fiber.Ctx) error { return controller.UpdateShopMenu(db, c) })
	app.Delete("/shopmenu/:id", func(c *fiber.Ctx) error { return controller.DeleteShopMenu(db, c) })

	//photo
	app.Post("/photos", func(c *fiber.Ctx) error { return controller.CreatePhoto(db, c) })
	app.Get("/photos/:id", func(c *fiber.Ctx) error { return controller.GetPhoto(db, c) })
	app.Get("/photos/menu/:menu_id", func(c *fiber.Ctx) error { return controller.GetPhotoByMenuID(db, c) })
	app.Get("/photos/shop/:shop_id", func(c *fiber.Ctx) error { return controller.GetPhotoByShopID(db, c) })
	app.Put("/photos/:id", func(c *fiber.Ctx) error { return controller.UpdatePhoto(db, c) })
	app.Delete("/photos/:id", func(c *fiber.Ctx) error { return controller.DeletePhoto(db, c) })

	//contact to admin
	app.Post("/contacts", func(c *fiber.Ctx) error { return controller.CreateContactToAdmin(db, c) })
	app.Get("/contacts/:id", func(c *fiber.Ctx) error { return controller.GetContactToAdmin(db, c) })
	app.Put("/contacts/:id", func(c *fiber.Ctx) error { return controller.UpdateContactToAdmin(db, c) })
	app.Delete("/contacts/:id", func(c *fiber.Ctx) error { return controller.DeleteContactToAdmin(db, c) })

	//filter
	//how to use filter-shops?keyword=coffee
	app.Get("/filter-shops", func(c *fiber.Ctx) error { return controller.FilterShopsByKeyword(db, c) })

	app.Get("/config", getENV)

	app.Listen(":8080")

}

// func ptrInt(i int) *int {
// 	return &i
// }

func getENV(c *fiber.Ctx) error {

	secret := os.Getenv("SECRET")
	if secret == "" {
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
		dsn = "user:12345678@tcp(152.42.240.7:3306)/BFM?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})

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
