package routes

import (
	"errors"
	"time"

	db "github.com/KompocikDot/cadeau/db/generated"
	"github.com/KompocikDot/cadeau/internal/server"
	"github.com/KompocikDot/cadeau/internal/types"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateAuthRoutes(sv *server.Server) {
	auth := sv.App.Group("/api/auth")
	auth.Post("/register/", func(c fiber.Ctx) error {
		p := new(types.UserRegister)

		if err := c.Bind().JSON(p); err != nil {
			return err
		}

		hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
		err := sv.DB.CreateUser(c.Context(), db.CreateUserParams{
			Username: p.Username,
			Password: string(hashedPwd),
		})

		if err != nil {
			return err
		}

		return nil
	})

	auth.Post("/login/", func(c fiber.Ctx) error {
		p := new(types.UserLogin)

		if err := c.Bind().JSON(p); err != nil {
			return err
		}

		u, err := sv.DB.GetUserByUsername(c, p.Username)
		if err != nil {
			return err
		}

		if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p.Password)) != nil {
			return errors.New("password does not match")
		}

		expiryDate := time.Now().Add(time.Hour * 72)
		claims := jwt.MapClaims{
			"id":  u.ID,
			"exp": expiryDate.Unix(),
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		c.Cookie(&fiber.Cookie{
			Name:     "token",
			Value:    t,
			Path:     "/",
			Domain:   "localhost",
			HTTPOnly: true,
			Expires:  expiryDate,
			SameSite: "strict",
		})

		return nil
	})

	auth.Get("/logout/", func(c fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:     "token",
			Value:    "",
			Path:     "/",
			Domain:   "localhost",
			HTTPOnly: true,
			Expires:  time.Time{},
			SameSite: "strict",
		})

		return nil
	})

}
