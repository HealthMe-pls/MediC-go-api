package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


func getPatients(db *gorm.DB, c *fiber.Ctx) error {
	var patients []Patient
	db.Find(&patients)
	return c.JSON(patients)
  }
  

  func getPatientID(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var patient Patient
	db.First(&patient, id)
	return c.JSON(patient)
  }

  func createPatient(db *gorm.DB, c *fiber.Ctx) error {
	patient := new(Patient)
	if err := c.BodyParser(patient); err != nil {
	  return err
	}
	db.Create(&patient)
	return c.JSON(patient)
  }

  func updatePatient(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	patient := new(Patient)
	db.First(&patient, id)
	if err := c.BodyParser(patient); err != nil {
	  return err
	}
	db.Save(&patient)
	return c.JSON(patient)
  }
  
  func deletePatient(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	db.Delete(&Patient{}, id)
	return c.SendString("Book successfully deleted")
  }

  func uploadImage(db *gorm.DB, c *fiber.Ctx) error {
	// Extract Patient ID from request params
	patientID := c.Params("id")

	// Validate that the patient exists
	var patient Patient
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
	image := Image{
		PatientID: patient.ID,
		ImagePath: ptrString(fileName),
	}
	if err := db.Create(&image).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save image record")
	}

	return c.SendString(file.Filename)
	}

func getPatientImages(db *gorm.DB, c *fiber.Ctx) error {
	patientID := c.Params("id")

	// Validate that the patient exists
	var patient Patient
	if err := db.First(&patient, patientID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Patient not found")
	}

	// Fetch all images associated with the patient, including related Patient data
	var images []Image
	if err := db.Preload("Patient").Where("patient_id = ?", patientID).Find(&images).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch images")
	}

	return c.JSON(images)
}
