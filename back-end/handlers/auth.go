package handlers

import (
	"database/sql"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"tivix-performance-tracker-backend/database"
	"tivix-performance-tracker-backend/middleware"
	"tivix-performance-tracker-backend/models"
	"tivix-performance-tracker-backend/utils"
)

var validate = validator.New()

func Register(c *fiber.Ctx) error {
	var req models.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Dados inválidos",
		})
	}

	if err := validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Dados de entrada inválidos",
			"details": err.Error(),
		})
	}

	var existingUser models.User
	err := database.DB.Get(&existingUser, "SELECT id FROM users WHERE email = $1", req.Email)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "Email já está em uso",
		})
	} else if err != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Erro interno do servidor",
		})
	}

	user := models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Name:      req.Name,
		Role:      "user", // Role padrão
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if req.Role != "" {
		user.Role = req.Role
	}

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
			"message": "Erro ao criar usuário",
		})
	}

	token, err := middleware.GenerateJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao gerar token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Usuário criado com sucesso",
		"data": models.LoginResponse{
			Token: token,
			User:  user,
		},
	})
}

func Login(c *fiber.Ctx) error {
	var req models.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Dados inválidos",
		})
	}

	if err := validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Dados de entrada inválidos",
			"details": err.Error(),
		})
	}

	var user models.User
	err := database.DB.Get(&user, "SELECT * FROM users WHERE email = $1", req.Email)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Credenciais inválidas",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Erro interno do servidor",
		})
	}

	if !user.IsActive {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   true,
			"message": "Usuário inativo",
		})
	}

	if err := user.CheckPassword(req.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Credenciais inválidas",
		})
	}

	token, err := middleware.GenerateJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao gerar token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Login realizado com sucesso",
		"data": models.LoginResponse{
			Token: token,
			User:  user,
		},
	})
}

func GetProfile(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(*middleware.JWTClaims)

	var user models.User
	err := database.DB.Get(&user, "SELECT * FROM users WHERE id = $1", userClaims.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Usuário não encontrado",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  user,
	})
}

func RefreshToken(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(*middleware.JWTClaims)

	var user models.User
	err := database.DB.Get(&user, "SELECT * FROM users WHERE id = $1", userClaims.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Usuário não encontrado",
		})
	}

	if !user.IsActive {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   true,
			"message": "Usuário inativo",
		})
	}

	token, err := middleware.GenerateJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao gerar token",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token": token,
		},
	})
}

func CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Dados inválidos",
		})
	}

	if err := validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Dados de entrada inválidos",
			"details": err.Error(),
		})
	}

	if err := utils.ValidatePassword(req.TemporaryPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	var existingUser models.User
	err := database.DB.Get(&existingUser, "SELECT id FROM users WHERE email = $1", req.Email)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "error",
			"message": "Email já está em uso",
		})
	} else if err != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro interno do servidor",
		})
	}

	user := models.User{
		ID:                  uuid.New(),
		Email:               req.Email,
		Name:                req.Name,
		Role:                req.Role,
		NeedsPasswordChange: true,
		IsActive:            true,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := user.HashPassword(req.TemporaryPassword); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao processar senha",
		})
	}

	query := `
		INSERT INTO users (id, email, password, name, role, needs_password_change, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = database.DB.Exec(query, user.ID, user.Email, user.Password, user.Name, user.Role, user.NeedsPasswordChange, user.IsActive, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao criar usuário",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

func SetNewPassword(c *fiber.Ctx) error {
	var req models.SetNewPasswordRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Dados inválidos",
		})
	}

	if err := validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Dados de entrada inválidos",
			"details": err.Error(),
		})
	}

	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	userClaims := c.Locals("user").(*middleware.JWTClaims)

	var user models.User
	err := database.DB.Get(&user, "SELECT * FROM users WHERE id = $1", userClaims.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Usuário não encontrado",
		})
	}

	if !user.NeedsPasswordChange {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Usuário não precisa trocar a senha",
		})
	}

	if err := user.HashPassword(req.NewPassword); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao processar nova senha",
		})
	}

	query := `
		UPDATE users 
		SET password = $1, needs_password_change = false, updated_at = $2 
		WHERE id = $3
	`
	_, err = database.DB.Exec(query, user.Password, time.Now(), user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao atualizar senha",
		})
	}

	user.NeedsPasswordChange = false
	user.UpdatedAt = time.Now()

	token, err := middleware.GenerateJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao gerar token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": models.LoginResponse{
			Token: token,
			User:  user,
		},
	})
}

func ChangePassword(c *fiber.Ctx) error {
	var req models.ChangePasswordRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Dados inválidos",
		})
	}

	if err := validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Dados de entrada inválidos",
			"details": err.Error(),
		})
	}

	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	userClaims := c.Locals("user").(*middleware.JWTClaims)

	var user models.User
	err := database.DB.Get(&user, "SELECT * FROM users WHERE id = $1", userClaims.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Usuário não encontrado",
		})
	}

	if err := user.CheckPassword(req.CurrentPassword); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Senha atual incorreta",
		})
	}

	if err := user.HashPassword(req.NewPassword); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao processar nova senha",
		})
	}

	query := `UPDATE users SET password = $1, updated_at = $2 WHERE id = $3`
	_, err = database.DB.Exec(query, user.Password, time.Now(), user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao atualizar senha",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Senha alterada com sucesso",
	})
}

func ListUsers(c *fiber.Ctx) error {
	var users []models.User

	query := `
		SELECT id, email, name, role, needs_password_change, is_active, created_at, updated_at 
		FROM users 
		ORDER BY created_at DESC
	`
	err := database.DB.Select(&users, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao buscar usuários",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   users,
	})
}
