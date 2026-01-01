package middleware

import (
	"errors"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
)

func CreateJWTMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
		Extractor:  extractors.FromCookie("token"),
		ErrorHandler: func(c fiber.Ctx, err error) error {
			if errors.Is(err, extractors.ErrNotFound) {
				return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired JWT")
			}
			if e, ok := err.(*fiber.Error); ok {
				return c.Status(e.Code).SendString(e.Message)
			}
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired JWT")
		},
	})
}
