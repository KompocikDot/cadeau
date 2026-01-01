package utils

import (
	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserId(c fiber.Ctx) int64 {
	user := jwtware.FromContext(c)
	claims := user.Claims.(jwt.MapClaims)
	return int64(claims["id"].(float64))
}
