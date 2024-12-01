package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID int `json: "id"`
	Name string `json: "name"`
	Email string `json: "email"`
	Password string `json: "password"`
}

var users []User

func getUsers(c *fiber.Ctx) error {
	return c.JSON(users)
}
func getUserID(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))

	if err != nil { 
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for _, user := range users{
		if userID == user.ID {
			return c.JSON(user)
		}
	} 
	return c.Status(fiber.StatusNotFound).SendString("user ID not found")
}

func createUser(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user) ; err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	user.ID = len(users) +1
	users = append(users, *user)
	return c.JSON(users) 
}

func updateUser(c *fiber.Ctx) error {
    userID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
    }

    userUpdate := new(User)
    if err := c.BodyParser(userUpdate); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
    }

    for i, user := range users {
        if userID == user.ID {
            if userUpdate.Name != "" {
                users[i].Name = userUpdate.Name
            }
            if userUpdate.Password != "" {
                users[i].Password = userUpdate.Password
            }
            if userUpdate.Email != "" {
                users[i].Email = userUpdate.Email
            }
            return c.JSON(users[i])
        }
    }
    return c.Status(fiber.StatusNotFound).SendString("User ID not found")
}

func deleteUser (c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))

	if err != nil { 
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for i, user := range users{
		if userID == user.ID {
			// delete the current index from slice
			users = append(users[:i], users[i+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}
	} 
	return c.Status(fiber.StatusNotFound).SendString("user ID not found")
}
