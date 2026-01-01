package routes

import (
	"database/sql"
	"errors"

	db "github.com/KompocikDot/cadeau/db/generated"
	"github.com/KompocikDot/cadeau/internal/middleware"
	"github.com/KompocikDot/cadeau/internal/server"
	"github.com/KompocikDot/cadeau/internal/types"
	"github.com/KompocikDot/cadeau/internal/utils"
	"github.com/gofiber/fiber/v3"
)

func CreateOccasionsRouter(sv *server.Server) {
	router := sv.App.Group("/occasions", middleware.CreateJWTMiddleware())

	router.Get("/", func(c fiber.Ctx) error {
		occasions, err := sv.DB.GetUserOccasionsByUserId(c, utils.GetUserId(c))
		if err != nil {
			return err
		}

		res := make([]types.OccasionResponse, len(occasions))
		for i, o := range occasions {
			res[i] = types.OccasionResponse{
				Name: o.Name,
				Id:   o.ID,
			}
		}

		return c.JSON(res)
	})

	router.Get("/:occasionId/", func(c fiber.Ctx) error {
		occasionId := fiber.Params[int64](c, "occasionId", 0)
		if occasionId <= 0 {
			return errors.New("invalid occasionId")
		}

		occasion, err := sv.DB.GetOccasionById(c, occasionId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.SendStatus(fiber.StatusNotFound)
			}

			return err
		}

		if occasion.GiftReceiver != utils.GetUserId(c) {
			return c.SendStatus(fiber.StatusForbidden)
		}

		guests, err := sv.DB.GetOccasionGuests(c, occasionId)
		if err != nil {
			return err
		}

		guestsRes := make([]types.UsersResponse, len(guests))
		for i, g := range guests {
			guestsRes[i] = types.UsersResponse{Id: g.ID, Username: g.Username}
		}

		return c.JSON(types.OccasionResponse{Name: occasion.Name, Id: occasion.ID, Guests: guestsRes})
	})

	router.Patch("/:occasionId/", func(c fiber.Ctx) error {
		occasionId := fiber.Params[int64](c, "occasionId", 0)
		if occasionId <= 0 {
			return errors.New("invalid occasionId")
		}

		u := new(types.UpdateOccasion)
		if err := c.Bind().JSON(u); err != nil {
			return err
		}

		err := sv.DB.UpdateOccasion(c, db.UpdateOccasionParams{
			Name: u.Name,
			ID:   occasionId,
		})

		if err != nil {
			return err
		}

		return c.JSON(types.OccasionResponse{Name: u.Name, Id: occasionId})
	})

	router.Post("/:occasionId/gifts/", func(c fiber.Ctx) error {
		g := new(types.CreateGift)
		if err := c.Bind().JSON(g); err != nil {
			return err
		}

		occasionId := fiber.Params[int64](c, "occasionId", 0)
		if occasionId <= 0 {
			return errors.New("invalid occasionId")
		}

		o, err := sv.DB.GetOccasionById(c, occasionId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.SendStatus(fiber.StatusNotFound)
			}

			return err
		}

		if o.GiftReceiver != utils.GetUserId(c) {
			return c.SendStatus(fiber.StatusForbidden)
		}

		newId, err := sv.DB.CreateGift(c, db.CreateGiftParams{
			Name:     g.Name,
			Url:      g.URL,
			Occasion: occasionId,
		})

		if err != nil {
			return err
		}

		return c.JSON(types.CreateGiftResponse{Id: newId})
	})

	router.Get("/:occasionId/gifts/", func(c fiber.Ctx) error {
		occasionId := fiber.Params[int64](c, "occasionId", 0)
		if occasionId <= 0 {
			return errors.New("invalid occasionId")
		}

		gifts, err := sv.DB.SelectGiftsByOcassionId(c, db.SelectGiftsByOcassionIdParams{
			Occasion:     occasionId,
			GiftReceiver: utils.GetUserId(c),
		})

		if err != nil {
			return err
		}

		res := make([]types.GiftResponse, len(gifts))
		for i, g := range gifts {
			res[i] = types.GiftResponse{Name: g.Name, URL: g.Url, Id: g.ID}
		}

		return c.JSON(res)
	})

	router.Delete("/:occasionId/gifts/:giftId/", func(c fiber.Ctx) error {
		giftId := fiber.Params[int64](c, "giftId", 0)
		if giftId <= 0 {
			return errors.New("invalid giftId")
		}

		g, err := sv.DB.GetGiftById(c, giftId)
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		if g.GiftReceiver.Int64 != utils.GetUserId(c) {
			return c.SendStatus(fiber.StatusForbidden)
		}

		if err := sv.DB.DeleteGift(c, giftId); err != nil {
			return err
		}

		return nil
	})

	router.Patch("/:occasionId/gifts/:giftId/", func(c fiber.Ctx) error {
		u := new(types.UpdateGift)
		if err := c.Bind().JSON(u); err != nil {
			return err
		}

		giftId := fiber.Params[int64](c, "giftId", 0)
		if giftId <= 0 {
			return errors.New("invalid giftId")
		}

		g, err := sv.DB.GetGiftById(c, giftId)
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		if g.GiftReceiver.Int64 != utils.GetUserId(c) {
			return c.SendStatus(fiber.StatusForbidden)
		}

		err = sv.DB.UpdateGift(c, db.UpdateGiftParams{
			ID:   giftId,
			Name: u.Name,
			Url:  u.URL,
		})
		if err != nil {
			return err
		}

		return nil
	})

	router.Post("/occasions/", func(c fiber.Ctx) error {
		p := new(types.CreateOccasion)
		if err := c.Bind().JSON(p); err != nil {
			return err
		}

		tx, err := sv.RawDB.BeginTx(c, nil)
		if err != nil {
			return err
		}

		defer tx.Rollback()

		newId, err := sv.DB.WithTx(tx).CreateOccasion(c, db.CreateOccasionParams{
			Name:         p.Name,
			GiftReceiver: utils.GetUserId(c),
		})

		if err != nil {
			return err
		}

		for _, guestId := range p.Guests {
			err := sv.DB.WithTx(tx).AssignUserToOccasion(c, db.AssignUserToOccasionParams{OccasionID: newId, UserID: guestId})
			if err != nil {
				return err
			}
		}

		tx.Commit()

		return c.JSON(types.CreateOccasionResponse{Id: newId})
	})

	router.Delete("/occasions/:occasionId/", func(c fiber.Ctx) error {
		occasionId := fiber.Params[int64](c, "occasionId", 0)
		if occasionId <= 0 {
			return errors.New("invalid occasionId")
		}

		o, err := sv.DB.GetOccasionById(c, occasionId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.SendStatus(fiber.StatusNotFound)
			}

			return err
		}

		if o.GiftReceiver != utils.GetUserId(c) {
			return c.SendStatus(fiber.StatusForbidden)
		}

		err = sv.DB.DeleteOccasion(c, occasionId)
		if err != nil {
			return err
		}

		return nil
	})
}
