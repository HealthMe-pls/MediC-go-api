package controller

import (
	"fmt"
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

func GetMarketMapDetail(db *gorm.DB, c *fiber.Ctx) error {
	// Retrieve all market maps
	var marketMaps []model.MarketMap
	if err := db.Find(&marketMaps).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve market maps")
	}

	// Iterate through market maps and add shop names
	var result []map[string]interface{}

	for _, marketMap := range marketMaps {
		// Retrieve the shop details using the ShopID
		var shop model.Shop
		if err := db.First(&shop, marketMap.ShopID).Error; err != nil {
			// If shop is not found, use "no shop"
			result = append(result, map[string]interface{}{
				"block_id":  marketMap.BlockID,
                "block_name": marketMap.BlockName,
                "block_zone": marketMap.BlockZone,
				"shop_id":   marketMap.ShopID,
				"shop_name": "no shop",
                "category_id": "no categoory",
			})
		} else {
			// Append the block_id and shop_name to the result
			result = append(result, map[string]interface{}{
				"block_id":  marketMap.BlockID,
                "block_name": marketMap.BlockName,
                "block_zone": marketMap.BlockZone,
				"shop_id":   marketMap.ShopID,
				"shop_name": shop.Name,
                "category_id": shop.ShopCategoryID,
			})
		}
	}

	// Return the result as a JSON response
	return c.JSON(result)
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
    // var shop model.Shop
    // if err := db.First(&shop, marketMap.ShopID).Error; err != nil {
    //     return c.Status(fiber.StatusNotFound).SendString("Shop not found")
    // }

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
    // Get the block_id from the URL parameter
    blockID := c.Params("id")

    // Convert block_id to integer (optional based on your application)
    blockIDUint, err := strconv.Atoi(blockID)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid block_id format",
        })
    }

    // Query the MarketMap to get the shop_id associated with the block_id
    var marketMap model.MarketMap
    if err := db.Where("block_id = ?", blockIDUint).First(&marketMap).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": "Market map not found for the provided block_id",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   "Failed to retrieve market map",
            "details": err.Error(),
        })
    }

    // Check if shop_id exists in the retrieved market map
    if marketMap.ShopID == nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "No shop associated with the provided block_id",
        })
    }

    // Call GetShopByID to fetch the shop details using the shop_id
    shopID := strconv.Itoa(int(*marketMap.ShopID)) // Convert uint to string
    c.Params("id", shopID) // Dynamically set the "id" parameter for the next handler

    return GetShopByID(db, c)
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
func UpdateAllMarketMaps(db *gorm.DB, c *fiber.Ctx) error {
    // Parse the request body for a list of updates
    var updates []fiber.Map
    if err := c.BodyParser(&updates); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Failed to parse request body")
    }

    for _, update := range updates {
        blockID, ok := update["block_id"].(float64) // Ensure block_id is provided and is a valid number
        if !ok {
            return c.Status(fiber.StatusBadRequest).SendString("Invalid or missing block_id in one of the updates")
        }

        // Convert block_id to uint
        blockIDUint := uint(blockID)

        // Fetch the existing MarketMap record
        var marketMap model.MarketMap
        if err := db.First(&marketMap, "block_id = ?", blockIDUint).Error; err != nil {
            return c.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("MarketMap with block_id %d not found", blockIDUint))
        }

        // Update fields dynamically from the update map
        for key, value := range update {
            if key == "block_id" {
                continue // Skip block_id to avoid accidental changes
            }

            // Dynamically set field values
            if err := db.Model(&marketMap).Update(key, value).Error; err != nil {
                return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to update field %s for block_id %d", key, blockIDUint))
            }
        }
    }

    // Return success response
    return c.Status(fiber.StatusOK).SendString("All MarketMaps updated successfully")
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

func DeleteMarketMapsByBlockName(db *gorm.DB, c *fiber.Ctx) error {
    // Get the BlockName from the URL parameters
    blockName := c.Params("block_name")

    // Delete all MarketMap entries for the specified BlockName
    if err := db.Where("block_name = ?", blockName).Delete(&model.MarketMap{}).Error; err != nil {
        // Return an internal server error if deletion fails
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete MarketMaps for BlockName")
    }

    // Successfully deleted all MarketMaps for the BlockName
    return c.SendString("All MarketMaps for the BlockName successfully deleted")
}

func UpdateMarketMapByBlockName(db *gorm.DB, c *fiber.Ctx) error {
    // Get the BlockName from the URL parameters
    blockName := c.Params("block_name")

    // Fetch the existing MarketMap record
    var marketMap model.MarketMap
    if err := db.First(&marketMap, "block_name = ?", blockName).Error; err != nil {
        return c.Status(fiber.StatusNotFound).SendString("MarketMap with specified BlockName not found")
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

func GetMapByBlockName(db *gorm.DB, c *fiber.Ctx) error {
    // Get the BlockName from the URL parameters
    blockName := c.Params("block_name")

    // Find the map by BlockName
    var marketMap model.MarketMap
    if err := db.Preload("Shop").First(&marketMap, "block_name = ?", blockName).Error; err != nil {
        if gorm.ErrRecordNotFound == err {
            return c.Status(fiber.StatusNotFound).SendString("Market map not found")
        }
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve market map")
    }

    // Return the market map as JSON
    return c.Status(fiber.StatusOK).JSON(marketMap)
}

