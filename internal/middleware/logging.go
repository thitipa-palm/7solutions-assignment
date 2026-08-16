package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logging(c *fiber.Ctx) error {
	start := time.Now()

	err := c.Next()

	duration := time.Since(start)

	log.Printf(
		"method=%s path=%s status=%d duration=%s",
		c.Method(),
		c.Path(),
		c.Response().StatusCode(),
		duration,
	)

	return err
}
