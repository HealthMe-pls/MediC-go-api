package controller

import (
	"strconv"
	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetMarketMap(db *gorm.DB,c *fiber.Ctx) error {
	var marketMap []model.MarketMap
	db.Find(&marketMap)
	return c.JSON(marketMap)
}

func CreateMarketMap(db *gorm.DB, c *fiber.Ctx) error {
    // Create an instance of the MarketMap struct to bind the incoming request
    marketMap := new(model.MarketMap)
    
    // Parse the incoming request body to bind it to the marketMap struct
    if err := c.BodyParser(marketMap); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Failed to parse request body",
        })
    }

    // Check if the Shop with the provided ShopID exists
    var shop model.Shop
    if err := db.First(&shop, marketMap.ShopID).Error; err != nil {
        return c.Status(fiber.StatusNotFound).SendString("Shop not found")
    }

    // Create the new MarketMap record in the database
    if result := db.Create(&marketMap); result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to create market map",
        })
    }

    // Return the newly created MarketMap as a JSON response
    return c.Status(fiber.StatusCreated).JSON(marketMap)
}

func GetShopInMapID(db *gorm.DB, c *fiber.Ctx) error {
    id := c.Params("id")  // Get the shop ID from the URL parameter
    var shop model.Shop

    // Query the Shop table by ID
    if err := db.Where("id = ?", id).First(&shop).Error; err != nil {
        // If no shop is found, return a 404 status with a message
        return c.Status(fiber.StatusNotFound).SendString("No shop found with this ID")
    }

    // Extract the shop name
    shopName := shop.Name

    // Return the shop name as JSON
    return c.JSON(fiber.Map{
        "shop_name": shopName,
    })
}

func DeleteMarketMapsByBlockID(db *gorm.DB, c *fiber.Ctx) error {
    // Get the BlockID from the URL (parameter)
    blockID := c.Params("block_id")

    // Convert BlockID to uint (ensure error handling)
    blockIDUint, err := strconv.Atoi(blockID)
    if err != nil {
        // Return a bad request error if conversion fails
        return c.Status(fiber.StatusBadRequest).SendString("Invalid BlockID format")
    }

    // Delete all MarketMap entries for the specific BlockID
    if err := db.Where("block_id = ?", blockIDUint).Delete(&model.MarketMap{}).Error; err != nil {
        // Return an internal server error if deletion fails
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete MarketMaps for BlockID")
    }

    // Successfully deleted all MarketMaps for the BlockID
    return c.SendString("All MarketMaps for the BlockID successfully deleted")
}

func UpdateMarketMapByBlockID(db *gorm.DB, c *fiber.Ctx) error {
    // Get the block_id from the URL parameters
    blockID := c.Params("block_id")

    // Convert block_id to uint
    blockIDUint, err := strconv.Atoi(blockID)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid BlockID format")
    }

    // Fetch the existing MarketMap record
    var marketMap model.MarketMap
    if err := db.First(&marketMap, "block_id = ?", blockIDUint).Error; err != nil {
        return c.Status(fiber.StatusNotFound).SendString("MarketMap with specified BlockID not found")
    }

    // Parse the request body to get updated fields
    if err := c.BodyParser(&marketMap); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Failed to parse request body")
    }

    // Save the updated MarketMap back to the database
    if err := db.Save(&marketMap).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to update MarketMap")
    }

    // Return the updated MarketMap as JSON
    return c.Status(fiber.StatusOK).JSON(marketMap)
}

func GetMapByBlockID(db *gorm.DB, c *fiber.Ctx) error {
    // Get the block ID from the URL parameters
    blockID := c.Params("id")

    // Convert the block ID to uint
    id, err := strconv.Atoi(blockID)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid Block ID format")
    }

    // Find the map by block ID
    var marketMap model.MarketMap
    if err := db.Preload("Shop").First(&marketMap, "block_id = ?", id).Error; err != nil {
        if gorm.ErrRecordNotFound == err {
            return c.Status(fiber.StatusNotFound).SendString("Market map not found")
        }
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve market map")
    }

    // Return the market map as JSON
    return c.Status(fiber.StatusOK).JSON(marketMap)
}

