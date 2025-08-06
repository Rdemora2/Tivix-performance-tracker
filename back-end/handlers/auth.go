package handlers

import (
	"database/sql"
	"fmt"
	"strings"
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
		INSERT INTO users (id, email, password, name, role, company_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = database.DB.Exec(query, user.ID, user.Email, user.Password, user.Name, user.Role, user.CompanyID, user.IsActive, user.CreatedAt, user.UpdatedAt)
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
	user := c.Locals("user").(*middleware.JWTClaims)
	
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

	// Verificar regras de empresa
	var finalCompanyID *uuid.UUID
	if user.Role == "admin" {
		// Admin pode especificar qualquer empresa (obrigatório agora)
		if req.CompanyID == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Admin deve especificar uma empresa para o usuário",
			})
		}
		finalCompanyID = req.CompanyID
	} else if user.Role == "manager" {
		// Manager só pode criar usuários na sua própria empresa
		if user.CompanyID == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Manager deve estar associado a uma empresa",
			})
		}
		finalCompanyID = user.CompanyID
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Apenas admins e managers podem criar usuários",
		})
	}

	// Verificar se a empresa existe
	var companyExists bool
	err := database.DB.Get(&companyExists, "SELECT EXISTS(SELECT 1 FROM companies WHERE id = $1 AND is_active = true)", *finalCompanyID)
	if err != nil || !companyExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Empresa não encontrada ou inativa",
		})
	}

	if err := utils.ValidatePassword(req.TemporaryPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	var existingUser models.User
	existingUserErr := database.DB.Get(&existingUser, "SELECT id FROM users WHERE email = $1", req.Email)
	if existingUserErr == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "error",
			"message": "Email já está em uso",
		})
	} else if existingUserErr != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro interno do servidor",
		})
	}

	newUser := models.User{
		ID:                  uuid.New(),
		Email:               req.Email,
		Name:                req.Name,
		Role:                req.Role,
		CompanyID:           finalCompanyID,
		NeedsPasswordChange: true,
		IsActive:            true,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := newUser.HashPassword(req.TemporaryPassword); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao processar senha",
		})
	}

	query := `
		INSERT INTO users (id, email, password, name, role, company_id, needs_password_change, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err = database.DB.Exec(query, newUser.ID, newUser.Email, newUser.Password, newUser.Name, newUser.Role, newUser.CompanyID, newUser.NeedsPasswordChange, newUser.IsActive, newUser.CreatedAt, newUser.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao criar usuário",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   newUser,
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
	user := c.Locals("user").(*middleware.JWTClaims)
	var users []models.User

	var query string
	var args []interface{}

	if user.Role == "admin" {
		// Admins podem ver todos os usuários
		query = `
			SELECT id, email, name, role, company_id, needs_password_change, is_active, created_at, updated_at 
			FROM users 
			ORDER BY created_at DESC
		`
	} else {
		// Managers e usuários só podem ver usuários da sua empresa
		if user.CompanyID == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Usuário deve estar associado a uma empresa",
			})
		}
		
		query = `
			SELECT id, email, name, role, company_id, needs_password_change, is_active, created_at, updated_at 
			FROM users 
			WHERE company_id = $1
			ORDER BY created_at DESC
		`
		args = append(args, *user.CompanyID)
	}

	err := database.DB.Select(&users, query, args...)
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

func UpdateUser(c *fiber.Ctx) error {
	currentUser := c.Locals("user").(*middleware.JWTClaims)
	id := c.Params("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ID do usuário inválido",
		})
	}

	var req models.UpdateUserRequest
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

	// Verificar se o usuário existe e se pode ser editado
	var existingUser models.User
	err = database.DB.Get(&existingUser, "SELECT * FROM users WHERE id = $1", userID)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Usuário não encontrado",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao buscar usuário",
		})
	}

	// Verificar permissões de edição
	if currentUser.Role != "admin" {
		// Managers só podem editar usuários da sua própria empresa
		if currentUser.CompanyID == nil || existingUser.CompanyID == nil || *currentUser.CompanyID != *existingUser.CompanyID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Sem permissão para editar este usuário",
			})
		}
		
		// Managers não podem editar admins
		if existingUser.Role == "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Managers não podem editar administradores",
			})
		}
		
		// Managers não podem promover usuários a admin
		if req.Role != nil && *req.Role == "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Managers não podem promover usuários a administrador",
			})
		}
	}

	// Construir query de update dinamicamente
	updates := []string{}
	args := []interface{}{}
	argCount := 1

	if req.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argCount))
		args = append(args, *req.Name)
		argCount++
	}

	if req.Email != nil {
		// Verificar se email já existe em outro usuário
		var emailExists bool
		err := database.DB.Get(&emailExists, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id != $2)", *req.Email, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Erro ao verificar email",
			})
		}
		if emailExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"status":  "error",
				"message": "Email já está em uso",
			})
		}

		updates = append(updates, fmt.Sprintf("email = $%d", argCount))
		args = append(args, *req.Email)
		argCount++
	}

	if req.Role != nil {
		updates = append(updates, fmt.Sprintf("role = $%d", argCount))
		args = append(args, *req.Role)
		argCount++
	}

	if req.CompanyID != nil {
		// Apenas admin pode mudar a empresa do usuário
		if currentUser.Role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Apenas administradores podem alterar a empresa do usuário",
			})
		}

		// Verificar se a empresa existe
		var companyExists bool
		err := database.DB.Get(&companyExists, "SELECT EXISTS(SELECT 1 FROM companies WHERE id = $1 AND is_active = true)", *req.CompanyID)
		if err != nil || !companyExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Empresa não encontrada ou inativa",
			})
		}

		updates = append(updates, fmt.Sprintf("company_id = $%d", argCount))
		args = append(args, *req.CompanyID)
		argCount++
	}

	if req.IsActive != nil {
		// Apenas admin pode ativar/desativar usuários
		if currentUser.Role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Apenas administradores podem ativar/desativar usuários",
			})
		}

		updates = append(updates, fmt.Sprintf("is_active = $%d", argCount))
		args = append(args, *req.IsActive)
		argCount++
	}

	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Nenhum campo foi fornecido para atualização",
		})
	}

	// Adicionar updated_at
	updates = append(updates, fmt.Sprintf("updated_at = $%d", argCount))
	args = append(args, time.Now())
	argCount++

	// Adicionar ID para WHERE clause
	args = append(args, userID)

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(updates, ", "), argCount)
	
	_, err = database.DB.Exec(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao atualizar usuário",
		})
	}

	// Buscar usuário atualizado
	var updatedUser models.User
	err = database.DB.Get(&updatedUser, "SELECT * FROM users WHERE id = $1", userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao buscar usuário atualizado",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   updatedUser,
	})
}

// DeleteUser exclui um usuário
func DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ID de usuário inválido",
		})
	}

	// Obter usuário atual das claims do JWT
	currentUser := c.Locals("user").(*middleware.JWTClaims)

	// Buscar o usuário a ser excluído
	var userToDelete models.User
	err = database.DB.Get(&userToDelete, "SELECT * FROM users WHERE id = $1", userUUID)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Usuário não encontrado",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao buscar usuário",
		})
	}

	// Verificar permissões de exclusão
	if currentUser.Role != "admin" {
		// Managers só podem excluir usuários da sua própria empresa
		if currentUser.CompanyID == nil || userToDelete.CompanyID == nil || *currentUser.CompanyID != *userToDelete.CompanyID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Sem permissão para excluir este usuário",
			})
		}
		
		// Managers não podem excluir admins ou outros managers
		if userToDelete.Role == "admin" || userToDelete.Role == "manager" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Sem permissão para excluir administradores ou gerentes",
			})
		}
	}

	// Verificar se o usuário está tentando excluir a si mesmo
	if currentUser.UserID == userUUID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Você não pode excluir sua própria conta",
		})
	}

	// Verificar se existem dados associados ao usuário (se necessário)
	// Por exemplo, verificar se o usuário criou algum relatório ou outro dado importante

	// Executar a exclusão
	_, err = database.DB.Exec("DELETE FROM users WHERE id = $1", userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao excluir usuário",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Usuário excluído com sucesso",
		"data": fiber.Map{
			"deletedUser": fiber.Map{
				"id":    userToDelete.ID,
				"name":  userToDelete.Name,
				"email": userToDelete.Email,
				"role":  userToDelete.Role,
			},
		},
	})
}
