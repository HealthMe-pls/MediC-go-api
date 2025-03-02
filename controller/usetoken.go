package controller

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)
func GetShopDetailsByLoggedInEntrepreneur(db *gorm.DB, c *fiber.Ctx) error {
    // Get the JWT token from the Authorization header
    tokenString := c.Get("Authorization")
    if tokenString == "" {
        fmt.Println("Missing token")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
    }

    // Remove "Bearer " prefix if present
    if len(tokenString) > len("Bearer ") && tokenString[:len("Bearer ")] == "Bearer " {
        tokenString = tokenString[len("Bearer "):]
    } else {
        fmt.Println("Invalid token format")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
    }

    // Parse and validate the token
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            fmt.Println("Unexpected signing method")
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte("your_secret_key"), nil
    })

    if err != nil || !token.Valid {
        fmt.Println("Invalid or expired token:", err)
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
    }

    // Extract entrepreneur name from token
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        fmt.Println("Invalid token claims")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
    }

    entrepreneurName, ok := claims["username"].(string)
    if !ok {
        fmt.Println("Entrepreneur name not found in token")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Entrepreneur name missing in token"})
    }

    // Log the extracted entrepreneur name
    fmt.Println("Entrepreneur Name: ", entrepreneurName)

    // Find entrepreneur by name
    var entrepreneur model.Entrepreneur
    if err := db.Where("username = ?", entrepreneurName).First(&entrepreneur).Error; err != nil {
        fmt.Println("Entrepreneur not found:", err)
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Entrepreneur not found"})
    }

    // Log found entrepreneur ID
    fmt.Println("Entrepreneur ID: ", entrepreneur.ID)

    // Fetch all shops belonging to this entrepreneur
    var shops []model.Shop
    if err := db.Preload("Entrepreneur").
        Preload("ShopCategory").
        Where("entrepreneur_id = ?", entrepreneur.ID).
        Find(&shops).Error; err != nil {
        fmt.Println("Failed to retrieve shops:", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve shops", "details": err.Error()})
    }

    // If no shops found, return an empty array
    if len(shops) == 0 {
        fmt.Println("No shops found")
        return c.JSON([]fiber.Map{})
    }

    // Prepare response array
    var shopResponses []fiber.Map
    for _, shop := range shops {
        shopOpenDates, _ := getShopOpenDates(db, shop.ID)
        shopMenus, _ := getShopMenus(db, shop.ID)
        socialMedias, _ := getSocialMedia(db, shop.ID)
        shopPhotos, _ := getPhotosByShopID(db, shop.ID)

        // Construct response object
        shopResponse := fiber.Map{
            "shop_id":         shop.ID,
            "name":            shop.Name,
            "entrepreneur_id": shop.Entrepreneur.ID,
            "category_id":     shop.ShopCategory.ID,
            "category":        shop.ShopCategory.Name,
            "open_status":     shop.OpenStatus,
            "description":     shop.Description,
            "photos":          shopPhotos,
            "shop_open_dates": shopOpenDates,
            "menus":           shopMenus,
            "social_media":    socialMedias,
        }
        shopResponses = append(shopResponses, shopResponse)
    }

    return c.JSON(shopResponses)
}


func GetShopDetailsByEntrepreneurID(db *gorm.DB, c *fiber.Ctx) error {
	// Parse Entrepreneur ID from request parameters
	entrepreneurID, err := strconv.ParseUint(c.Params("entrepreneur_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid entrepreneur ID",
		})
	}

	// Fetch all shops belonging to the entrepreneur
	var shops []model.Shop
	if err := db.Preload("Entrepreneur").
		Preload("ShopCategory").
		Where("entrepreneur_id = ?", entrepreneurID).
		Find(&shops).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve shops",
			"details": err.Error(),
		})
	}

	// If no shops found, return an empty array
	if len(shops) == 0 {
		return c.JSON([]fiber.Map{})
	}

	// Prepare response array
	var shopResponses []fiber.Map

	// Iterate over each shop to fetch related data
	for _, shop := range shops {
		shopOpenDates, err := getShopOpenDates(db, shop.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to retrieve shop open dates",
				"details": err.Error(),
			})
		}

		shopMenus, err := getShopMenus(db, shop.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to retrieve shop menus",
				"details": err.Error(),
			})
		}

		socialMedias, err := getSocialMedia(db, shop.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to retrieve social media",
				"details": err.Error(),
			})
		}

		shopPhotos, err := getPhotosByShopID(db, shop.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to retrieve shop photos",
				"details": err.Error(),
			})
		}

		// Construct individual shop response
		shopResponse := fiber.Map{
			"shop_id":         shop.ID,
			"name":            shop.Name,
			"entrepreneur_id": shop.Entrepreneur.ID,
			"category_id":     shop.ShopCategory.ID,
			"category":        shop.ShopCategory.Name,
			"open_status":     shop.OpenStatus,
			"description":     shop.Description,
			"photos":          shopPhotos,
			"shop_open_dates": shopOpenDates,
			"menus":           shopMenus,
			"social_media":    socialMedias,
		}

		// Append to response array
		shopResponses = append(shopResponses, shopResponse)
	}

	return c.JSON(shopResponses)
}