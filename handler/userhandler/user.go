package userhandler

import (
	"BackendCoursyclopedia/model/usermodel"
	usersvc "BackendCoursyclopedia/service/userservice"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type IUserHandler interface {
	GetUsers(c *fiber.Ctx) error
	GetOneUser(c *fiber.Ctx) error
	CreateOneUser(c *fiber.Ctx) error
	DeleteOneUser(c *fiber.Ctx) error
	UpdateOneUser(c *fiber.Ctx) error
	GetUserByEmail(c *fiber.Ctx) error
	DropAllUsers(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type UserHandler struct {
	UserService usersvc.IUserService
}

func NewUserHandler(userService usersvc.IUserService) IUserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (h UserHandler) withTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout()
	defer cancel()

	users, err := h.UserService.GetAllUsers(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	for i := range users {
		users[i].Password = ""
	}

	return c.JSON(fiber.Map{
		"message": "Users retrieved successfully",
		"data":    users,
	})
}

func (h *UserHandler) GetOneUser(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout()
	defer cancel()

	userID := c.Params("id") // Assuming the user ID is passed as a URL parameter
	user, err := h.UserService.GetUserByID(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	user.Password = ""

	return c.JSON(fiber.Map{
		"message": "Specific User retrieved successfully",
		"data":    user,
	})
}

func (h *UserHandler) GetUserByEmail(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout()
	defer cancel()

	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}

	user, err := h.UserService.GetUserByEmail(ctx, email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	user.Password = ""

	return c.JSON(fiber.Map{
		"message": "User retrieved successfully",
		"data":    user,
	})
}

func (h *UserHandler) CreateOneUser(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout()
	defer cancel()

	var user usermodel.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	createdUser, err := h.UserService.CreateNewUser(ctx, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// To ensure the password hash doesn't get sent back, reset it to an empty string
	createdUser.Password = ""

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"data":    createdUser,
	})
}

func (h *UserHandler) DeleteOneUser(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout()
	defer cancel()

	userID := c.Params("id") // Retrieve the userID from the URL parameter.
	err := h.UserService.DeleteSpecificUser(ctx, userID)
	if err != nil {
		// If an error occurred, send an appropriate response.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// If no error, send a success response.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

func (h *UserHandler) UpdateOneUser(c *fiber.Ctx) error {
	// Context with timeout for the operation
	ctx, cancel := h.withTimeout()
	defer cancel()

	// Extract the user ID from the URL parameter
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is required"})
	}

	// Parse the JSON body into a User struct
	var updateUser usermodel.User
	if err := c.BodyParser(&updateUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Could not parse request body"})
	}

	// Call the UserService to update the user
	updatedUser, err := h.UserService.UpdateSpecificByID(ctx, userID, updateUser)
	if err != nil {
		// Handle specific errors like "user not found" or "invalid input" differently if needed
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	updatedUser.Password = ""

	// Return the updated user and a success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
		"data":    updatedUser,
	})
}

func (h *UserHandler) DropAllUsers(c *fiber.Ctx) error {
	ctx, cancel := h.withTimeout()
	defer cancel()

	if err := h.UserService.DropAllUsers(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the login request body
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	// Use the UserService.Login method to authenticate the user and generate a JWT token
	user, token, err := h.UserService.Login(c.Context(), loginRequest.Email, loginRequest.Password)
	if err != nil {
		// For security, don't reveal whether the email or password was incorrect
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Login successful, return the JWT token in the response
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"data":    user.Profile,
		"token":   token, // Include the JWT token in the response
	})
}

// func (h *UserHandler) UpdateOneUser(c *fiber.Ctx) error {
// 	// Context with timeout for the operation
// 	ctx, cancel := h.withTimeout()
// 	defer cancel()

// 	// Extract the user ID from the URL parameter
// 	userID := c.Params("id")
// 	if userID == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is required"})
// 	}

// 	// Parse the JSON body into a User struct
// 	var updateUser model.User
// 	if err := c.BodyParser(&updateUser); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Could not parse request body"})
// 	}

// 	// Call the UserService to update the user
// 	updatedUser, err := h.UserService.UpdateSpecificByID(ctx, userID, updateUser)
// 	if err != nil {
// 		// Handle specific errors like "user not found" or "invalid input" differently if needed
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	// Return the updated user and a success message
// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "User updated successfully",
// 		"data":    updatedUser,
// 	})
// }
