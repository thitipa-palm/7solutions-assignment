package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/thitipa-palm/7solutions-assignment/internal/handler"
)

func Setup(
	app *fiber.App,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	authMiddleware fiber.Handler,
) {
	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	users := api.Group("/users", authMiddleware) //NOTE - apply authMiddleware สำหรับทุก route ของ /users
	users.Post("/", userHandler.Create)
	users.Get("/:id", userHandler.GetByID)
	users.Get("/", userHandler.List)
	users.Patch("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)
}
