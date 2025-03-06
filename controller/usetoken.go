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
    // Get TempShops using the helper function
    tempShops, err := GetTempShopsByEntrepreneurID(db, entrepreneur.ID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"temp_shops": tempShops})
}

// GetTempShopsByEntrepreneurID retrieves TempShops for all shops owned by an entrepreneur
func GetTempShopsByEntrepreneurID(db *gorm.DB, entrepreneurID uint) ([]fiber.Map, error) {
	var shops []model.Shop
	if err := db.Where("entrepreneur_id = ?", entrepreneurID).Find(&shops).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve shops: %w", err)
	}

	if len(shops) == 0 {
		return nil, fmt.Errorf("no shops found for entrepreneur_id: %d", entrepreneurID)
	}

	var tempShops []model.TempShop
	var tempShopResponses []fiber.Map

	// Collect Shop IDs
	var shopIDs []uint
	for _, shop := range shops {
		shopIDs = append(shopIDs, shop.ID)
	}

	// Fetch all TempShops for the given shop IDs
	if err := db.Where("shop_id IN (?)", shopIDs).Find(&tempShops).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve temp shops: %w", err)
	}

	for _, tempShop := range tempShops {
		if tempShop.ShopID == nil {
			continue
		}
		shopID := *tempShop.ShopID

		// Fetch related data
		TempSocials, _ := GetTempSocialsByShopID(db, shopID)
		tempMenus, _ := getMenuForTempByShopID(db, shopID)
		photoData, _ := getPhotoForTempByShopID(db, shopID)
		tempTimes, _ := GetTempTimeByShopID(db, shopID)
		shopOpenDates, _ := GetShopOpenDateForTempByShopID(db, shopID)

		// Extract data from maps
		photosShop, _ := photoData["photos_shop"].([]fiber.Map)
		photosMenu, _ := photoData["photos_menu"].([]fiber.Map)

		socials, _ := TempSocials["socials"].([]fiber.Map)
		addTime, _ := tempTimes["addTime"].([]fiber.Map)
		editTime, _ := tempTimes["editTime"].([]fiber.Map)
		deleteTime, _ := tempTimes["deleteTime"].([]fiber.Map)

		// Construct the response
		tempShopResponses = append(tempShopResponses, fiber.Map{
			"id":            tempShop.TempID,
            "entrepreneur_id":entrepreneurID,
			"name":          tempShop.Name,
			"description":	 tempShop.Description,
			"category_id":	 tempShop.ShopCategoryID,
			"shop_id":       shopID,
			"socials":       socials,
			"menus":         tempMenus,
			"photos_shop":   photosShop,
			"photos_menu":   photosMenu,
			"addTime":       addTime,
			"editTime":      editTime,
			"deleteTime":    deleteTime,
			"time":          shopOpenDates,
		})
	}

	return tempShopResponses, nil
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
func GetEntrepreneurByIDLogin(db *gorm.DB, c *fiber.Ctx) error {
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

    // Extract entrepreneur ID from token
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        fmt.Println("Invalid token claims")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
    }

    entrepreneurName, ok := claims["username"].(string)
    if !ok {
        fmt.Println("Entrepreneur name not found in token")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Entrepreneur ID missing in token"})
    }

    // Query the database for the entrepreneur with the extracted ID
    var entrepreneur model.Entrepreneur
    if err := db.Where("username = ?", entrepreneurName).First(&entrepreneur).Error; err != nil {
        fmt.Println("Entrepreneur not found:", err)
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Entrepreneur not found"})
    }
	decryptedPassword, err := DecryptPassword(entrepreneur.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to decrypt password",
				"details": err.Error(),
			})
		}
	entrepreneur.Password = decryptedPassword

    // If successful, return the entrepreneur data as a JSON response
    return c.JSON(entrepreneur)
}