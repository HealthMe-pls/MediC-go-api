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
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Missing token",
        })
    }
    
    // Remove the "Bearer " prefix
    if len(tokenString) > len("Bearer ") && tokenString[:len("Bearer ")] == "Bearer " {
        tokenString = tokenString[len("Bearer "):]
    } else {
        fmt.Println("Invalid token format")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid token format",
        })
    }

    // Parse and validate the token
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Ensure the token is signed with HMAC algorithm (or your preferred algorithm)
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            fmt.Println("Unexpected signing method")
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        // Return the secret key used for signing the token
        return []byte("your_secret_key"), nil
    })
    
    // Token parsing and validation
    if err != nil {
        fmt.Println("Error parsing token:", err)
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid or expired token",
        })
    }
    if !token.Valid {
        fmt.Println("Token is invalid or expired")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid or expired token",
        })
    }

    // Extract the EntrepreneurID from the token claims
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        fmt.Println("Invalid token claims")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid token claims",
        })
    }

    entrepreneurID, ok := claims["entrepreneur_id"].(float64) // type assertion to float64
    if !ok {
        fmt.Println("Entrepreneur ID not found in token")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Entrepreneur ID missing in token",
        })
    }

    // Log the extracted entrepreneur ID
    fmt.Println("Entrepreneur ID: ", entrepreneurID)

    // Fetch all shops belonging to the logged-in entrepreneur
    var shops []model.Shop
    if err := db.Preload("Entrepreneur").
        Preload("ShopCategory").
        Where("entrepreneur_id = ?", uint(entrepreneurID)).
        Find(&shops).Error; err != nil {
        fmt.Println("Failed to retrieve shops:", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   "Failed to retrieve shops",
            "details": err.Error(),
        })
    }

    // If no shops found, return an empty array
    if len(shops) == 0 {
        fmt.Println("No shops found")
        return c.JSON([]fiber.Map{})
    }

    // Prepare the response array
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

        // Construct the individual shop response
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

        // Append to the response array
        shopResponses = append(shopResponses, shopResponse)
    }

    // Return the shop details as a response
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