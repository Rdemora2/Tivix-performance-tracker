package main

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"tivix-performance-tracker-backend/database"
	"tivix-performance-tracker-backend/routes"
)

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Conectar ao banco de dados
	database.Connect()

	// Executar migrações
	database.Migrate()

	// Criar instância do Fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	// Configurar CORS com variáveis de ambiente
	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "http://localhost:5173"
	}

	// Lista de origens permitidas (incluindo desenvolvimento e produção)
	allowedOrigins := []string{
		corsOrigin,
		"http://localhost:3000",
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"https://performancetracker.tivix.com.br",
		"https://performance.valiantgroup.com.br",
	}

	// Remover duplicatas e vazios
	var finalOrigins []string
	seen := make(map[string]bool)
	for _, origin := range allowedOrigins {
		if origin != "" && !seen[origin] {
			finalOrigins = append(finalOrigins, origin)
			seen[origin] = true
		}
	}

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(finalOrigins, ","),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Rotas
	routes.SetupRoutes(app)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Tivix Performance Tracker API is running",
		})
	})

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
