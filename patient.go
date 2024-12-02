package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Patient struct {
	ID *int `json:"ID"`
	Name *string `json:"Name"`
	Email *string `json:"Email"`
	Age *int `json:"Age"`
}

var patients []Patient

func getPatients(c *fiber.Ctx) error {
	return c.JSON(patients)
}
func getPatientID(c *fiber.Ctx) error {
	patientID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for _, patient := range patients {
		if patient.ID != nil && *patient.ID == patientID { // Check for nil and compare values
			return c.JSON(patient)
		}
	}
	return c.Status(fiber.StatusNotFound).SendString("Patient ID not found")
}


func createPatient(c *fiber.Ctx) error {
	patient := new(Patient)
	if err := c.BodyParser(patient); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	newID := len(patients) + 1
	patient.ID = &newID // Assign a pointer to the new ID
	patients = append(patients, *patient)
	return c.JSON(patients)
}


func updatePatient(c *fiber.Ctx) error {
    patientID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid patient ID")
    }

    patientUpdate := new(Patient)
    if err := c.BodyParser(patientUpdate); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
    }

    for i, patient := range patients {
        if patientID == *patient.ID {
            if patientUpdate.Name != nil {
                patients[i].Name = patientUpdate.Name
            }
            if patientUpdate.Age != nil {
                patients[i].Age = patientUpdate.Age
            }
            if patientUpdate.Email != nil {
                patients[i].Email = patientUpdate.Email
            }
            return c.JSON(patients[i])
        }
    }
    return c.Status(fiber.StatusNotFound).SendString("Patient ID not found")
}

func deletePatient (c *fiber.Ctx) error {
	patientID, err := strconv.Atoi(c.Params("id"))

	if err != nil { 
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for i, patient := range patients{
		if patientID == *patient.ID {
			// delete the current index from slice
			patients = append(patients[:i], patients[i+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}
	} 
	return c.Status(fiber.StatusNotFound).SendString("patient ID not found")
}


func uploadImage(c *fiber.Ctx) error {
	// Read file from request
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Save the file to the server
	err = c.SaveFile(file, "./uploads/" + file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString(file.Filename)
}