package handler

import (
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/thitipa-palm/7solutions-assignment/internal/repository"
	"github.com/thitipa-palm/7solutions-assignment/internal/service"
)

type UserHandler struct {
	userService *service.UserService
	validate    *validator.Validate
}

func NewUserHandler(
	userService *service.UserService,
	validate *validator.Validate,
) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validate,
	}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var request CreateUserRequest

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
				Error: "invalid user data",
			},
		)
	}

	user, err := h.userService.Create(
		c.UserContext(),
		service.CreateUserInput{
			Name:     request.Name,
			Email:    request.Email,
			Password: request.Password,
		},
	)
	if err != nil {
		return handleUserError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	user, err := h.userService.GetByID(
		c.UserContext(),
		c.Params("id"),
	)
	if err != nil {
		return handleUserError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	users, err := h.userService.List(c.UserContext())
	if err != nil {
		return handleUserError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	var request UpdateUserRequest

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
				Error: "invalid user data",
			},
		)
	}

	user, err := h.userService.Update(
		c.UserContext(),
		c.Params("id"),
		service.UpdateUserInput{
			Name:  request.Name,
			Email: request.Email,
		},
	)
	if err != nil {
		return handleUserError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	err := h.userService.Delete(
		c.UserContext(),
		c.Params("id"),
	)
	if err != nil {
		return handleUserError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func handleUserError(
	c *fiber.Ctx,
	err error,
) error {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return c.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Error: "invalid user data",
			},
		)

	case errors.Is(err, service.ErrInvalidUserID):
		return c.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Error: "invalid user ID",
			},
		)

	case errors.Is(err, repository.ErrUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(
			ErrorResponse{
				Error: "user not found",
			},
		)

	case errors.Is(err, repository.ErrEmailAlreadyExists):
		return c.Status(fiber.StatusConflict).JSON(
			ErrorResponse{
				Error: "email already exists",
			},
		)

	default:
		log.Printf("handle user request: %v", err)

		return c.Status(
			fiber.StatusInternalServerError,
		).JSON(
			ErrorResponse{
				Error: "internal server error",
			},
		)
	}
}
