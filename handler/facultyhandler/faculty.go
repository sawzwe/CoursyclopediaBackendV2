package facultyhandler

import (
	"BackendCoursyclopedia/model/facultymodel"
	"encoding/json"
	"io"

	facultysvc "BackendCoursyclopedia/service/facultyservice"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type IFacultyHandler interface {
	GetFaculties(c *fiber.Ctx) error
	GetEachFaculty(c *fiber.Ctx) error
	GetMajorsForeachFaculty(c *fiber.Ctx) error
	CreateFaculty(c *fiber.Ctx) error
	UpdateFaculty(c *fiber.Ctx) error
	DeleteFaculty(c *fiber.Ctx) error
}

type FacultyHandler struct {
	FacultyService facultysvc.IFacultyService
}

func NewFacultyHandler(facultyService facultysvc.IFacultyService) IFacultyHandler {
	return &FacultyHandler{
		FacultyService: facultyService,
	}
}

func (h FacultyHandler) withTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

func (h FacultyHandler) GetFaculties(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout()
	defer cancel()

	faculties, err := h.FacultyService.GetAllFaculties(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Faculties retrieved successfully",
		"data":    faculties,
	})
}

func (h *FacultyHandler) GetEachFaculty(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout()
	defer cancel()

	facultyID := c.Params("id")
	faculty, err := h.FacultyService.GetFacultyByID(ctx, facultyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Specific Faculty retrieved successfully",
		"data":    faculty,
	})
}

func (h *FacultyHandler) GetMajorsForeachFaculty(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout()
	defer cancel()

	facultyID := c.Params("id")
	majors, err := h.FacultyService.GetMajorsForFaculty(ctx, facultyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Majors related to the faculty retrieved successfully",
		"data":    majors,
	})
}

func (h *FacultyHandler) CreateFaculty(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout()
	defer cancel()

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Image upload error"})
	}

	fileData, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process image"})
	}
	defer fileData.Close()

	imageBytes, err := io.ReadAll(fileData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read image data"})
	}

	facultyName := c.FormValue("FacultyName")
	if facultyName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Faculty name is required"})
	}

	faculty := facultymodel.Faculty{
		FacultyName: facultyName,
		Image:       imageBytes,
	}

	createdFaculty, err := h.FacultyService.CreateFaculty(ctx, faculty, imageBytes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Faculty created successfully",
		"data":    createdFaculty,
	})
}

func (h *FacultyHandler) UpdateFaculty(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout()
	defer cancel()

	facultyID := c.Params("id")
	var faculty facultymodel.Faculty
	if err := c.BodyParser(&faculty); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Image upload error"})
	}

	fileData, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process image"})
	}
	defer fileData.Close()

	imageBytes, err := io.ReadAll(fileData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read image data"})
	}

	facultyData := c.FormValue("faculty")
	if err := json.Unmarshal([]byte(facultyData), &faculty); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid faculty data"})
	}

	updatedFaculty, err := h.FacultyService.UpdateFaculty(ctx, facultyID, faculty, imageBytes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Faculty updated successfully",
		"data":    updatedFaculty,
	})
}

func (h FacultyHandler) DeleteFaculty(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout()
	defer cancel()

	facultyID := c.Params("id")
	err := h.FacultyService.DeleteFaculty(ctx, facultyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Faculty deleted successfully",
	})
}
