package controller

import (
	"log"
	"github.com/gofiber/fiber/v2"
    "golang.org/x/oauth2"
	"context"
	"os"
)

var infomaniakConfig = &oauth2.Config{
    ClientID:     os.Getenv("INFOMANIAK_CLIENT_ID"),
    ClientSecret: os.Getenv("INFOMANIAK_CLIENT_SECRET"),
    RedirectURL:  "http://localhost:80/auth/callback", // Frontend callback
    Endpoint: oauth2.Endpoint{
        AuthURL:  "https://login.infomaniak.com/oauth2/authorize",
        TokenURL: "https://login.infomaniak.com/oauth2/token",
    },
    Scopes:       []string{"email", "profile"},
}


func loginInfomaniak(c *fiber.Ctx) error {
    url := infomaniakConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
    return c.Redirect(url)
}

func callbackInfomaniak(c *fiber.Ctx) error {
    code := c.Query("code")
    if code == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Code not found"})
    }

    token, err := infomaniakConfig.Exchange(context.Background(), code)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get token"})
    }

    client := infomaniakConfig.Client(context.Background(), token)
    resp, err := client.Get("https://api.infomaniak.com/userinfo")
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user info"})
    }

    defer resp.Body.Close()
    return c.JSON(fiber.Map{"token": token.AccessToken})
}

func main() {
    app := fiber.New()

    app.Get("/auth/login", loginInfomaniak)
    app.Get("/auth/callback", callbackInfomaniak)

    log.Fatal(app.Listen(":8080"))
}