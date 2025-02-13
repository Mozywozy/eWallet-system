package server

import (
	"ewallet-engine/internal/auth"
	"ewallet-engine/internal/balance"
	"ewallet-engine/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)

	userRepo := auth.NewUserRepository(s.db)
	authService := auth.NewAuthService(userRepo)
	authHandler := auth.NewAuthHandler(authService)

	// Routing
	api := s.App.Group("/user/v1")
	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)
	api.Post("/logout", auth.JWTMiddleware() ,authHandler.Logout)
	api.Post("/refresh", authHandler.RefreshToken)

}

func (s *FiberServer) BalanceFiberRoutes() {
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	db := database.New().GetDB()
	balanceRepo := balance.NewBalanceRepository(db)
	balanceService := balance.NewBalanceService(balanceRepo)
	balanceHandler := balance.NewBalanceHandler(balanceService)

	api := s.App.Group("/user/v1")
	api.Get("/balance", auth.JWTMiddleware(), balanceHandler.GetBalanceHandler)
	api.Post("/topup", auth.JWTMiddleware(), balanceHandler.TopUpBalanceHandler)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
