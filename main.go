package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// use godotenv to get .env variables
	if err := godotenv.Load(); err!= nil { // gogotenv init
		log.Fatal("load .env error")
	}

	// ตัวแทนการสื่อสารกับ http server
	app := fiber.New() // fiber init

	// Static file serving to access uploaded images via /upload/{filename}
	app.Static("/upload", "./uploads")

	// Apply CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Adjust this to be more restrictive if needed
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// for test without database server
	patients = append(patients, Patient{
		ID:    ptrInt(1),
		Name:  ptrString("John Doe"),
		Email: ptrString("johndoe@gmail.com"),
		Age:   ptrInt(60),
	})
	patients = append(patients, Patient{
		ID:    ptrInt(2),
		Name:  ptrString("Lorem Ipsum"),
		Email: ptrString("loremipsum@gmail.com"),
		Age:   ptrInt(49),
	})
						// c = response and request fiber context
	app.Get("/hello" , func(c *fiber.Ctx) error {return c.SendString("test gogo")})

	// test route 
	app.Get("/patient", getPatients) 
	app.Get("/patient/:id", getPatientID)
	app.Post("/patient", createPatient)
	app.Put("/patient/:id", updatePatient)
	app.Delete("/patient/:id", deletePatient)
	app.Post("/upload", uploadImage)
	app.Get("/config", getENV)

	app.Listen(":8080")
}

func ptrString(s string) *string {
	return &s
}

func ptrInt(i int) *int {
	return &i
}

func getENV (c *fiber.Ctx) error {
	
	secret := os.Getenv("SECRET")
	if secret == ""{
		secret = "defaultsecret"
	}
	return c.JSON(fiber.Map{ 
		"SECRET": secret,
	})

}