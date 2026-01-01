package server

import (
	"database/sql"
	"log"
	"net/url"

	db "github.com/KompocikDot/cadeau/db/generated"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

type Server struct {
	App   *fiber.App
	RawDB *sql.DB
	DB    *db.Queries
}

func CreateServer() *Server {
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

	return &Server{
		App:   app,
		RawDB: conn,
		DB:    db.New(conn),
	}
}
