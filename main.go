package main

import (
	"log"

	"github.com/KompocikDot/cadeau/internal/middleware"
	"github.com/KompocikDot/cadeau/internal/routes"
	"github.com/KompocikDot/cadeau/internal/server"
	"github.com/KompocikDot/cadeau/internal/types"
	"github.com/KompocikDot/cadeau/internal/utils"
	"github.com/gofiber/fiber/v3"

	_ "github.com/amacneil/dbmate/v2/pkg/driver/sqlite"
	_ "modernc.org/sqlite"
)

func main() {
	sv := server.CreateServer()
	jwt := middleware.CreateJWTMiddleware()

	routes.CreateAuthRoutes(sv)
	routes.CreateOccasionsRouter(sv)

	usersApi := sv.App.Group("/api/users").Use(jwt)
	usersApi.Get("/", func(c fiber.Ctx) error {
		users, err := sv.DB.GetUsersList(c, utils.GetUserId(c))
		if err != nil {
			return err
		}

		res := make([]types.UsersResponse, len(users))
		for i, u := range users {
			res[i] = types.UsersResponse{Id: u.ID, Username: u.Username}
		}

		return c.JSON(res)
	})

	userApi := sv.App.Group("/app/user/me", jwt)
	userApi.Get("/", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	if err := sv.App.Listen(":8000"); err != nil {
		log.Fatalln(err)
	}
}
