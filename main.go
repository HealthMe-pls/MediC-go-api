package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// ตัวแทนการสื่อสารกับ http server
	app := fiber.New() // fiber init

	// Apply CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Adjust this to be more restrictive if needed
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	users = append(users, User{ID: 1,Name: "John Doe", Email: "johndoe@gmail.com", Password: "1234"})
	users = append(users, User{ID: 2,Name: "Lorem Ipsum", Email: "loremipsum@gmail.com", Password: "5678"})
						// c = response and request fiber context
	app.Get("/hello" , func(c *fiber.Ctx) error {return c.SendString("test gogo")})

	app.Get("/user", getUsers) 
	app.Get("/user/:id", getUserID)
	app.Post("/user", createUser)
	app.Put("/user/:id", updateUser)
	app.Delete("/user/:id", deleteUser)

	app.Listen(":8080")
}