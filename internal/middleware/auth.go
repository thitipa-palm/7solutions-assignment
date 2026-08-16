package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/thitipa-palm/7solutions-assignment/internal/service"
)

const AuthenticatedUserIDKey = "userID"

func Authenticate(
	tokenService *service.TokenService,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Get(fiber.HeaderAuthorization)

		parts := strings.Fields(authorization)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") {
			return unauthorized(c)
		}

		tokenString := parts[1]

		claims, err := tokenService.Parse(tokenString)
		if err != nil {
			return unauthorized(c)
		}

		c.Locals(
			AuthenticatedUserIDKey,
			claims.Subject, //NOTE - userId
		)

		return c.Next()
	}
}

func unauthorized(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(
		fiber.Map{
			"error": "unauthorized",
		},
	)
}
