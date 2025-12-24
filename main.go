package main

import (
	"database/sql"
	"errors"
	"log"
	"net/url"

	db "github.com/KompocikDot/cadeau/db/generated"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	"github.com/gofiber/fiber/v3"

	_ "github.com/amacneil/dbmate/v2/pkg/driver/sqlite"
	_ "modernc.org/sqlite"
)

type User struct {
	Username string
	Password string `json:"-"`
}
type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		AppName:       "CadeauAPI",
	})

	u, _ := url.Parse("sqlite:db.sqlite3")
	migrate := dbmate.New(u)

	err := migrate.CreateAndMigrate()
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := sql.Open("sqlite", "file:db.sqlite3")
	if err != nil {
		log.Fatalln(err)
	}

	dbC := db.New(conn)

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/auth/register", func(c fiber.Ctx) error {
		p := new(UserRegister)

		if err := c.Bind().JSON(p); err != nil {
			return err
		}

		err := dbC.CreateUser(c.Context(), db.CreateUserParams{
			Username: p.Username,
			Password: p.Password,
		})

		if err != nil {
			return err
		}

		return nil
	})

	app.Post("/auth/login", func(c fiber.Ctx) error {
		p := new(UserLogin)

		if err := c.Bind().JSON(p); err != nil {
			return err
		}

		u, err := dbC.GetUser(c.Context(), p.Username)
		if err != nil {
			return err
		}

		if u.Password != p.Password {
			return errors.New("password does not match")
		}

		return nil
	})

	app.Listen(":3000")
}
