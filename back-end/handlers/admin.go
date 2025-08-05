package handlers

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"tivix-performance-tracker-backend/database"
	"tivix-performance-tracker-backend/models"
)

// CreateAdminUser cria o primeiro usuário administrador (apenas se não houver usuários no sistema)
func CreateAdminUser(c *fiber.Ctx) error {
	// Verificar se já existem usuários no sistema
	var userCount int
	err := database.DB.Get(&userCount, "SELECT COUNT(*) FROM users")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao verificar usuários existentes",
		})
	}

	if userCount > 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   true,
			"message": "Sistema já possui usuários cadastrados",
		})
	}

	// Verificar se a chave de instalação está correta
	type InitRequest struct {
		InstallKey string `json:"installKey"`
		Email      string `json:"email" validate:"required,email"`
		Password   string `json:"password" validate:"required,min=6"`
		Name       string `json:"name" validate:"required,min=2"`
	}

	var req InitRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Dados inválidos",
		})
	}

	// Validar chave de instalação
	expectedKey := os.Getenv("INSTALL_KEY")
	if expectedKey == "" {
		expectedKey = "TIVIX_INSTALL_2024"
	}

	if req.InstallKey != expectedKey {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Chave de instalação inválida",
		})
	}

	// Validar dados
	if err := validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Dados de entrada inválidos",
			"details": err.Error(),
		})
	}

	// Criar usuário admin
	user := models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Name:      req.Name,
		Role:      "admin",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Hash da senha
	if err := user.HashPassword(req.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao processar senha",
		})
	}

	query := `
		INSERT INTO users (id, email, password, name, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = database.DB.Exec(query, user.ID, user.Email, user.Password, user.Name, user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao criar usuário administrador",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Usuário administrador criado com sucesso",
		"data": fiber.Map{
			"userId": user.ID,
			"email":  user.Email,
			"name":   user.Name,
			"role":   user.Role,
		},
	})
}

// CheckInitialization verifica se o sistema já foi inicializado
func CheckInitialization(c *fiber.Ctx) error {
	var userCount int
	err := database.DB.Get(&userCount, "SELECT COUNT(*) FROM users")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao verificar inicialização",
		})
	}

	return c.JSON(fiber.Map{
		"error":       false,
		"initialized": userCount > 0,
		"userCount":   userCount,
	})
}
