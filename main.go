package main

import (
	"database/sql"
	"errors"
	"log"
	"net/url"

	db "github.com/KompocikDot/cadeau/db/generated"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"

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

type CreateOccasion struct {
	Name         string `json:"name"`
	GiftReceiver int64  `json:"giftReceiver"`
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

	app.Post("/auth/register/", func(c fiber.Ctx) error {
		p := new(UserRegister)

		if err := c.Bind().JSON(p); err != nil {
			return err
		}

		hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
		err := dbC.CreateUser(c.Context(), db.CreateUserParams{
			Username: p.Username,
			Password: string(hashedPwd),
		})

		if err != nil {
			return err
		}

		return nil
	})

	app.Post("/auth/login/", func(c fiber.Ctx) error {
		p := new(UserLogin)

		if err := c.Bind().JSON(p); err != nil {
			return err
		}

		u, err := dbC.GetUserByUsername(c.Context(), p.Username)
		if err != nil {
			return err
		}

		if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p.Password)) != nil {
			return errors.New("password does not match")
		}

		return nil
	})

	app.Post("/occasions/", func(c fiber.Ctx) error {
		p := new(CreateOccasion)

		if err := c.Bind().JSON(p); err != nil {
			return err
		}

		err := dbC.CreateOccasion(c.Context(), db.CreateOccasionParams{
			Name:         p.Name,
			GiftReceiver: p.GiftReceiver,
		})
		if err != nil {
			return err
		}

		return nil
	})

	app.Get("/users/:userId/occasions/", func(c fiber.Ctx) error {
		userId := fiber.Params[int64](c, "userId", 0)
		if userId == 0 {
			return errors.New("invalid userId")
		}

		occasions, err := dbC.GetUserOccasionsByUserId(c.Context(), userId)
		if err != nil {
			return err
		}

		return c.JSON(occasions)
	})

	app.Get("/occasions/:occasionId/", func(c fiber.Ctx) error {
		userId := fiber.Params[int64](c, "userId", 0)
		if userId == 0 {
			return errors.New("invalid userId")
		}

		occasionId := fiber.Params[int64](c, "occasionId", 0)
		if occasionId == 0 {
			return errors.New("invalid occasionId")
		}

		occasion, err := dbC.GetOccasionById(c.Context(), db.GetOccasionByIdParams{
			ID:           occasionId,
			GiftReceiver: userId,
		})
		if err != nil {
			return err
		}

		return c.JSON(occasion)
	})

	app.Post("/gifts", func(c fiber.Ctx) error { return nil })

	app.Listen(":3000")
}
