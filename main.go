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
		&model.TempShop{},
		&model.ShopOpenDate{},
		&model.TempShopOpenDate{},
		&model.MarketMap{},
		&model.SocialMedia{},
		&model.ShopMenu{},
		&model.Workshop{},
		&model.Photo{},
		&model.TempMenu{},
		&model.TempSocial{},
		&model.DeletePhoto{},
		&model.DeleteSocial{},
		&model.DeleteMenu{}); err != nil {
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
	app.Delete("/entrepreneur/:id", func(c *fiber.Ctx) error { return controller.DeleteEntrepreneurByID(db, c) })

	//map --check
	app.Get("/map", func(c *fiber.Ctx) error { return controller.GetMarketMap(db, c) })
	app.Get("/mapdetail", func(c *fiber.Ctx) error { return controller.GetMarketMapDetail(db, c) })
	app.Get("/map/:id", func(c *fiber.Ctx) error { return controller.GetMapByBlockID(db, c) })
	app.Get("/shopInmap/:id", func(c *fiber.Ctx) error { return controller.GetShopInMapID(db, c) })
	app.Post("/map", func(c *fiber.Ctx) error { return controller.CreateMarketMap(db, c) })
	app.Delete("/map/:block_id", func(c *fiber.Ctx) error { return controller.DeleteMarketMapsByBlockID(db, c) })
	app.Put("/map/:block_id", func(c *fiber.Ctx) error { return controller.UpdateMarketMapByBlockID(db, c) })
	app.Put("/Allmap", func(c *fiber.Ctx) error { return controller.UpdateAllMarketMaps(db, c) })

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
	app.Get("/shopdetail", func(c *fiber.Ctx) error { return controller.GetShopDetail(db, c) })
	app.Get("/shopdetail/:id", func(c *fiber.Ctx) error { return controller.GetShopDetailByID(db, c) })
	app.Get("/entrepreneur/shopdetail/:entrepreneur_id", func(c *fiber.Ctx) error { return controller.GetShopDetailsByEntrepreneurID(db, c) })

	app.Get("/shop", func(c *fiber.Ctx) error { return controller.GetShops(db, c) })
	app.Put("/admin/shop/:id", func(c *fiber.Ctx) error { return controller.UpdateShopByAdmin(db, c) })
	app.Put("/shop/:shop_id", func(c *fiber.Ctx) error { return controller.UpdateTempShopByShopID(db, c) })
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
	app.Post("/social", func(c *fiber.Ctx) error { return controller.CreateSocialMediaByAdmin(db, c) })
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
	app.Post("/shopmenu", func(c *fiber.Ctx) error { return controller.CreateShopMenuByAdmin(db, c) })
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
	// photo create by entrepreneur
	app.Post("/photosmenu/:menu_id", func(c *fiber.Ctx) error {return controller.CreatePhotoByMenuID(db, c)})
	app.Post("/photosshop/:shop_id", func(c *fiber.Ctx) error {return controller.CreatePhotoByShopID(db, c)})
	//photo create by admin
	app.Post("/photos/workshop/:workshop_id", func(c *fiber.Ctx) error {return controller.CreatePhotoByWorkshopID(db, c)})
	app.Post("/photos/menu/:menu_id", func(c *fiber.Ctx) error {return controller.AdminCreatePhotoByMenuID(db, c)})
	app.Post("/photos/shop/:shop_id", func(c *fiber.Ctx) error {return controller.AdminCreatePhotoByShopID(db, c)})
	


	//contact to admin
	app.Post("/contacts/:entrepreneur_id", func(c *fiber.Ctx) error { return controller.CreateContactToAdmin(db, c) })
	app.Get("/contacts", func(c *fiber.Ctx) error { return controller.GetAllContacts(db, c) })
	// app.Post("/contacts", func(c *fiber.Ctx) error { return controller.CreateContactToAdmin(db, c) })
	app.Get("/contacts/:id", func(c *fiber.Ctx) error { return controller.GetContactToAdmin(db, c) })
	app.Put("/contacts/:id", func(c *fiber.Ctx) error { return controller.UpdateContactToAdmin(db, c) })
	app.Delete("/contacts/:id", func(c *fiber.Ctx) error { return controller.DeleteContactToAdmin(db, c) })

	//nothing usefull
	// TempShop routes
	app.Post("/tempshops", func(c *fiber.Ctx) error { return controller.CreateTempShop(db, c) })
	app.Get("/tempshops", func(c *fiber.Ctx) error { return controller.GetTempShops(db, c) })
	app.Get("/tempshops/:id", func(c *fiber.Ctx) error { return controller.GetTempShopByID(db, c) })
	app.Put("/tempshops/:id", func(c *fiber.Ctx) error { return controller.UpdateTempShop(db, c) })
	app.Delete("/tempshops/:id", func(c *fiber.Ctx) error { return controller.DeleteTempShop(db, c) })
	// TempShopOpenDate routes
	//use this 
	app.Post("/tempshopopendates", func(c *fiber.Ctx) error { return controller.CreateTempShopOpenDate(db, c) })
	app.Get("/tempshopopendates", func(c *fiber.Ctx) error { return controller.GetTempShopOpenDates(db, c) })
	app.Get("/tempshopopendates/:id", func(c *fiber.Ctx) error { return controller.GetTempShopOpenDateByID(db, c) })
	app.Put("/tempshopopendates/:id", func(c *fiber.Ctx) error { return controller.UpdateTempShopOpenDate(db, c) })
	app.Delete("/tempshopopendates/:id", func(c *fiber.Ctx) error { return controller.DeleteTempShopOpenDate(db, c) })

	//create with temp
	app.Post("/socials/entrepreneur", func(c *fiber.Ctx) error {return controller.CreateSocialWithTemp(db, c, false) })
	app.Post("/menus/entrepreneur", func(c *fiber.Ctx) error {return controller.CreateMenuWithTemp(db, c, false) })
	//admin
	app.Post("/socials/admin", func(c *fiber.Ctx) error {return controller.CreateSocialWithTemp(db, c, true)})
	app.Post("/menus/admin", func(c *fiber.Ctx) error {return controller.CreateMenuWithTemp(db, c, true) })
	app.Post("/createshop", func(c *fiber.Ctx) error {return controller.CreateShopWithTemp(db, c)})
	//menuupdate by entrepreneur
	app.Put("/updatemenu/:menu_id", func(c *fiber.Ctx) error {return controller.UpdateTempMenuByMenuID(db, c)})
	//social update by entrepreneur
	app.Put("updatesocial/:social_id",func(c *fiber.Ctx) error {return controller.UpdateSocialBySocialID(db, c)})
	//Bin
	// Bin routes
	//entrepreneur delete menu
	app.Post("/menubin", func(c *fiber.Ctx) error { return controller.CreateBinMenu(db, c) })
	app.Get("/menubin", func(c *fiber.Ctx) error { return controller.GetBinMenus(db, c) })
	app.Get("/menubin/:id", func(c *fiber.Ctx) error { return controller.GetBinMenuByID(db, c) })
	app.Get("/menubin/temp/:temp_id", func(c *fiber.Ctx) error { return controller.GetBinMenuByTempID(db, c) })
	// app.Put("/menubin/:id", func(c *fiber.Ctx) error { return controller.UpdateBinMenu(db, c) })
	app.Delete("/menubin/:id", func(c *fiber.Ctx) error { return controller.DeleteBinMenu(db, c) })
	app.Delete("/menubin/temp/:temp_id", func(c *fiber.Ctx) error { return controller.DeleteBinMenuByTempID(db, c) }) // 🔥 Delete by TempID

	//entrepreneur delete photo
	app.Post("/photobin", func(c *fiber.Ctx) error { return controller.CreateBinPhoto(db, c) })
	//nothing usefull this for test
	app.Get("/photobin", func(c *fiber.Ctx) error { return controller.GetBinPhotos(db, c) })
	app.Get("/photobin/:id", func(c *fiber.Ctx) error { return controller.GetBinPhotoByID(db, c) })
	app.Get("/photobin/temp/:temp_id", func(c *fiber.Ctx) error { return controller.GetBinPhotoByTempID(db, c) })
	// app.Put("/photobin/:id", func(c *fiber.Ctx) error { return controller.UpdateBinPhoto(db, c) })
	app.Delete("/photobin/:id", func(c *fiber.Ctx) error { return controller.DeleteBinPhoto(db, c) })
	app.Delete("/photobin/temp/:temp_id", func(c *fiber.Ctx) error { return controller.DeleteBinPhotoByTempID(db, c) }) // 🔥 Delete by TempID
	
	//entrepreneur delete social
	app.Post("/socialbin", func(c *fiber.Ctx) error { return controller.CreateBinSocial(db, c) })
	app.Get("/socialbin", func(c *fiber.Ctx) error { return controller.GetBinSocials(db, c) })
	app.Get("/socialbin/:id", func(c *fiber.Ctx) error { return controller.GetBinSocialByID(db, c) })
	app.Get("/socialbin/temp/:temp_id", func(c *fiber.Ctx) error { return controller.GetBinSocialByTempID(db, c) })
	// app.Put("/socialbin/:id", func(c *fiber.Ctx) error { return controller.UpdateBinSocial(db, c) })
	app.Delete("/socialbin/:id", func(c *fiber.Ctx) error { return controller.DeleteBinSocial(db, c) })
	app.Delete("/socialbin/temp/:temp_id", func(c *fiber.Ctx) error { return controller.DeleteBinSocialByTempID(db, c) }) // 🔥 Delete by TempID
	
	//not available
	app.Get("/availablemenus", func(c *fiber.Ctx) error { return controller.GetAvailableMenus(db, c) })
	app.Get("/availablemenus/:shop_id", func(c *fiber.Ctx) error { return controller.GetAvailableMenusByShopID(db, c) })
	app.Get("/availablephotos/menu/:menu_id", func(c *fiber.Ctx) error { return controller.GetAvailablePhotosByMenuID(db, c) })
	app.Get("/availablephotos/shop/:shop_id", func(c *fiber.Ctx) error { return controller.GetAvailablePhotosByShopID(db, c) })
	app.Get("/availablesocial/:shop_id", func(c *fiber.Ctx) error { return controller.GetAvailableSocialByShopID(db, c) })
	app.Get("/availableshopDetail/:shop_id", func(c *fiber.Ctx) error { return controller.GetAvailableShopDetailByID(db, c) })


	//admin manage
	app.Get("/approve/:id", func(c *fiber.Ctx) error { return controller.Handleapprove(db, c) })
	app.Put("/notApprove/:temp_id", func(c *fiber.Ctx) error { return controller.HandleNotApprove(db, c) })

	//filter
	//how to use search-shops?keyword=coffee
	app.Get("/search-shops", func(c *fiber.Ctx) error { return controller.SearchShopsByKeyword(db, c) })

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
		dsn = "user:12345678@tcp(127.0.0.1:3306)/BFM?charset=utf8mb4&parseTime=True&loc=Local"
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
