package main

import (
	"database/sql"
	"errors"
	"log"
	"net/url"
	"time"

	db "github.com/KompocikDot/cadeau/db/generated"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
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
	Name string `json:"name"`
}

func main() {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		AppName:       "CadeauAPI",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,
	}))

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

	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
		Extractor:  extractors.FromCookie("token"),
	})

	api := app.Group("/api")
	userApi := api.Group("/user").Use(jwtMiddleware)

	userApi.Get("/me/occasions/", func(c fiber.Ctx) error {
		user := jwtware.FromContext(c)
		claims := user.Claims.(jwt.MapClaims)
		userId := int64(claims["id"].(float64))

		occasions, err := dbC.GetUserOccasionsByUserId(c.Context(), userId)
		if err != nil {
			return err
		}

		return c.JSON(occasions)
	})

	userApi.Get("/me/occasions/:occasionId/", func(c fiber.Ctx) error {
		user := jwtware.FromContext(c)
		claims := user.Claims.(jwt.MapClaims)
		userId := int64(claims["id"].(float64))

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

	userApi.Post("/me/occasions/", func(c fiber.Ctx) error {
		p := new(CreateOccasion)
		if err := c.Bind().JSON(p); err != nil {
			return err
		}

		user := jwtware.FromContext(c)
		claims := user.Claims.(jwt.MapClaims)
		userId := int64(claims["id"].(float64))

		err := dbC.CreateOccasion(c.Context(), db.CreateOccasionParams{
			Name:         p.Name,
			GiftReceiver: userId,
		})
		if err != nil {
			return err
		}

		return nil
	})

	auth := api.Group("/auth")
	auth.Post("/register/", func(c fiber.Ctx) error {
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

	auth.Post("/login/", func(c fiber.Ctx) error {
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

		claims := jwt.MapClaims{
			"id":  u.ID,
			"exp": time.Now().Add(time.Hour * 72).Unix(),
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
		})

		return nil
	})

	if err := app.Listen(":8000"); err != nil {
		log.Fatalln(err)
	}
}
