package server

import (
	"hroost/server/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func (s *Server) start() error {
	s.app = fiber.New(fiber.Config{
		AppName: s.cfg.appName,
	})

	s.app.Use(logger.New())
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		// AllowOrigins: "http://localhost:3000,https://google.com",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// check where the request is coming from, then translate it into an application ID
	s.app.Use(middleware.AppId)

	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("oke")
	})

	presentation, err := s.newAppProvider()
	if err != nil {
		return err
	}

	s.presentation = presentation
	s.route()

	return nil
}
