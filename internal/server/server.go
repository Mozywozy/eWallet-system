package server

import (
	"github.com/gofiber/fiber/v2"

	"ewallet-engine/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "ewallet-engine",
			AppName:      "ewallet-engine",
		}),

		db: database.New(),
	}

	return server
}
