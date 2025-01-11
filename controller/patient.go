package controller

import (
	"fmt"

	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetPatients(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var patients []model.Patient
		db.Find(&patients)
		return c.JSON(patients)
	}
}

func GetPatientID(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var patient model.Patient
		db.First(&patient, id)
		return c.JSON(patient)
	}
}

func CreatePatient(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		patient := new(model.Patient)
		if err := c.BodyParser(patient); err != nil {
			return err
		}
		db.Create(&patient)
		return c.JSON(patient)
	}
}

func UpdatePatient(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		patient := new(model.Patient)
		db.First(&patient, id)
		if err := c.BodyParser(patient); err != nil {
			return err
		}
		db.Save(&patient)
		return c.JSON(patient)
	}
}

func DeletePatient(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		db.Delete(&model.Patient{}, id)
		return c.SendString("Book successfully deleted")
	}
}

func UploadImage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract Patient ID from request params
		patientID := c.Params("id")

		// Validate that the patient exists
		var patient model.Patient
		if err := db.First(&patient, patientID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Patient not found")
		}

		// Read file from request
		file, err := c.FormFile("image")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Failed to read image")
		}

		// Save the file to the server
		filePath := fmt.Sprintf("./uploads/%s", file.Filename)
		fileName := file.Filename
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save image")
		}

		// Save image record in the database
		image := model.Image{
			PatientID: patient.ID,
			ImagePath: ptrString(fileName),
		}
		if err := db.Create(&image).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save image record")
		}

		return c.SendString(file.Filename)
	}
}

func GetPatientImages(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		patientID := c.Params("id")

		// Validate that the patient exists
		var patient model.Patient
		if err := db.First(&patient, patientID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Patient not found")
		}

		// Fetch all images associated with the patient, including related Patient data
		var images []model.Image
		if err := db.Preload("Patient").Where("patient_id = ?", patientID).Find(&images).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch images")
		}

		return c.JSON(images)
	}
}
