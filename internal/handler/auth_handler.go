package handler

import (
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/thitipa-palm/7solutions-assignment/internal/repository"
	"github.com/thitipa-palm/7solutions-assignment/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(
	authService *service.AuthService,
	validate *validator.Validate,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validate,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var request RegisterRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Error: "invalid request body",
			},
		)
	}

	if err := h.validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Error: "invalid registration data",
			},
		)
	}

	user, err := h.authService.Register(
		c.UserContext(),
		service.RegisterInput{
			Name:     request.Name,
			Email:    request.Email,
			Password: request.Password,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			return c.Status(fiber.StatusBadRequest).JSON(
				ErrorResponse{
					Error: "invalid registration data",
				},
			)

		case errors.Is(err, repository.ErrEmailAlreadyExists):
			return c.Status(fiber.StatusConflict).JSON(
				ErrorResponse{
					Error: "email already exists",
				},
			)

		default:
			log.Printf("register user: %v", err)

			return c.Status(
				fiber.StatusInternalServerError,
			).JSON(
				ErrorResponse{
					Error: "internal server error",
				},
			)
		}
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var request LoginRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Error: "invalid request body",
			},
		)
	}

	if err := h.validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Error: "invalid login data",
			},
		)
	}

	token, err := h.authService.Login(
		c.UserContext(),
		service.LoginInput{
			Email:    request.Email,
			Password: request.Password,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			return c.Status(fiber.StatusUnauthorized).JSON(
				ErrorResponse{
					Error: "invalid email or password",
				},
			)

		default:
			log.Printf("login user: %v", err)

			return c.Status(
				fiber.StatusInternalServerError,
			).JSON(
				ErrorResponse{
					Error: "internal server error",
				},
			)
		}
	}

	return c.Status(fiber.StatusOK).JSON(
		LoginResponse{
			AccessToken: token,
			TokenType:   "Bearer",
		},
	)
}
