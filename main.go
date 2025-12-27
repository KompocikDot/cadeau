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
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
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

type OccasionResponse struct {
	Name string `json:"name"`
	Id   int64  `json:"id"`
}

type CreateOccasionResponse struct {
	Id int64 `json:"id"`
}

type UpdateOccasion struct {
	Name string `json:"name"`
}

type CreateGift struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type GiftResponse struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Id   int64  `json:"id"`
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
	}), helmet.New(), logger.New())

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

	api := app.Group("/api")
	userApi := api.Group("/user").Use(jwtMiddleware)

	userApi.Get("/me/", jwtMiddleware, func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	userApi.Get("/me/occasions/", func(c fiber.Ctx) error {
		user := jwtware.FromContext(c)
		claims := user.Claims.(jwt.MapClaims)
		userId := int64(claims["id"].(float64))

		occasions, err := dbC.GetUserOccasionsByUserId(c, userId)
		if err != nil {
			return err
		}

		res := make([]OccasionResponse, len(occasions))
		for i, o := range occasions {
			res[i] = OccasionResponse{
				Name: o.Name,
				Id:   o.ID,
			}
		}

		return c.JSON(res)
	})

	userApi.Get("/me/occasions/:occasionId/", func(c fiber.Ctx) error {
		user := jwtware.FromContext(c)
		claims := user.Claims.(jwt.MapClaims)
		userId := int64(claims["id"].(float64))

		occasionId := fiber.Params[int64](c, "occasionId", 0)
		if occasionId <= 0 {
			return errors.New("invalid occasionId")
		}

		occasion, err := dbC.GetOccasionById(c, db.GetOccasionByIdParams{
			ID:           occasionId,
			GiftReceiver: userId,
		})
		if err != nil {
			return err
		}

		return c.JSON(OccasionResponse{Name: occasion.Name, Id: occasion.ID})
	})

	userApi.Patch("/me/occasions/:occasionId/", func(c fiber.Ctx) error {
		occasionId := fiber.Params[int64](c, "occasionId", 0)
		if occasionId <= 0 {
			return errors.New("invalid occasionId")
		}

		u := new(UpdateOccasion)
		if err := c.Bind().JSON(u); err != nil {
			return err
		}

		err := dbC.UpdateOccasion(c, db.UpdateOccasionParams{
			Name: u.Name,
			ID:   occasionId,
		})

		if err != nil {
			return err
		}

		return c.JSON(OccasionResponse{Name: u.Name, Id: occasionId})
	})

	userApi.Post("/me/occasions/:occasionId/gifts/", func(c fiber.Ctx) error {
		g := new(CreateGift)
		if err := c.Bind().JSON(g); err != nil {
			return err
		}

		occasionId := fiber.Params[int64](c, "occasionId", 0)
		if occasionId <= 0 {
			return errors.New("invalid occasionId")
		}

		user := jwtware.FromContext(c)
		claims := user.Claims.(jwt.MapClaims)
		userId := int64(claims["id"].(float64))

		_, err := dbC.GetOccasionById(c, db.GetOccasionByIdParams{ID: occasionId, GiftReceiver: userId})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.SendStatus(fiber.StatusForbidden)
			}

			return err
		}

		err = dbC.CreateGift(c, db.CreateGiftParams{
			Name:     g.Name,
			Url:      g.URL,
			Occasion: occasionId,
		})

		if err != nil {
			return err
		}

		return nil
	})

	userApi.Get("/me/occasions/:occasionId/gifts/", func(c fiber.Ctx) error {
		occasionId := fiber.Params[int64](c, "occasionId", 0)
		if occasionId <= 0 {
			return errors.New("invalid occasionId")
		}

		user := jwtware.FromContext(c)
		claims := user.Claims.(jwt.MapClaims)
		userId := int64(claims["id"].(float64))

		gifts, err := dbC.SelectGiftsByOcassionId(c, db.SelectGiftsByOcassionIdParams{
			Occasion:     occasionId,
			GiftReceiver: userId,
		})

		if err != nil {
			return err
		}

		res := make([]GiftResponse, len(gifts))
		for i, g := range gifts {
			res[i] = GiftResponse{Name: g.Name, URL: g.Url, Id: g.ID}
		}

		return c.JSON(res)
	})

	userApi.Delete("/me/occasions/:occasionId/gifts/:giftId/", func(c fiber.Ctx) error {
		giftId := fiber.Params[int64](c, "giftId", 0)
		if giftId <= 0 {
			return errors.New("invalid giftId")
		}

		g, err := dbC.GetGiftById(c, giftId)
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		user := jwtware.FromContext(c)
		claims := user.Claims.(jwt.MapClaims)
		userId := int64(claims["id"].(float64))

		if g.GiftReceiver.Int64 != userId {
			return c.SendStatus(fiber.StatusForbidden)
		}

		if err := dbC.DeleteGift(c, giftId); err != nil {
			return err
		}

		return nil
	})

	userApi.Post("/me/occasions/", func(c fiber.Ctx) error {
		p := new(CreateOccasion)
		if err := c.Bind().JSON(p); err != nil {
			return err
		}

		user := jwtware.FromContext(c)
		claims := user.Claims.(jwt.MapClaims)
		userId := int64(claims["id"].(float64))

		newId, err := dbC.CreateOccasion(c, db.CreateOccasionParams{
			Name:         p.Name,
			GiftReceiver: userId,
		})
		if err != nil {
			return err
		}

		return c.JSON(CreateOccasionResponse{Id: newId})
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

		u, err := dbC.GetUserByUsername(c, p.Username)
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

	if err := app.Listen(":8000"); err != nil {
		log.Fatalln(err)
	}
}
