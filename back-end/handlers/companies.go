package handlers

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"tivix-performance-tracker-backend/database"
	"tivix-performance-tracker-backend/middleware"
	"tivix-performance-tracker-backend/models"
)

func CreateCompany(c *fiber.Ctx) error {
	var req models.CreateCompanyRequest

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

	// Verificar se já existe uma empresa com o mesmo nome
	var existingCompany models.Company
	err := database.DB.Get(&existingCompany, "SELECT id FROM companies WHERE name = $1", req.Name)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "error",
			"message": "Já existe uma empresa com esse nome",
		})
	} else if err != sql.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro interno do servidor",
		})
	}

	company := models.Company{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `
		INSERT INTO companies (id, name, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = database.DB.Exec(query, company.ID, company.Name, company.Description, company.IsActive, company.CreatedAt, company.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao criar empresa",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   company,
	})
}

func GetAllCompanies(c *fiber.Ctx) error {
	user := c.Locals("user").(*middleware.JWTClaims)
	var companies []models.Company

	var query string
	var args []interface{}

	if user.Role == "admin" {
		// Admins podem ver todas as empresas
		query = `
			SELECT id, name, description, is_active, created_at, updated_at
			FROM companies
			ORDER BY name ASC
		`
	} else {
		// Managers e usuários só podem ver sua própria empresa
		if user.CompanyID == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Usuário deve estar associado a uma empresa",
			})
		}
		
		query = `
			SELECT id, name, description, is_active, created_at, updated_at
			FROM companies
			WHERE id = $1
			ORDER BY name ASC
		`
		args = append(args, *user.CompanyID)
	}

	err := database.DB.Select(&companies, query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao buscar empresas",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   companies,
	})
}

func GetCompanyByID(c *fiber.Ctx) error {
	id := c.Params("id")
	companyID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ID da empresa inválido",
		})
	}

	var company models.Company
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM companies
		WHERE id = $1
	`
	err = database.DB.Get(&company, query, companyID)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Empresa não encontrada",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao buscar empresa",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   company,
	})
}

func UpdateCompany(c *fiber.Ctx) error {
	id := c.Params("id")
	companyID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ID da empresa inválido",
		})
	}

	var req models.UpdateCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Dados inválidos",
		})
	}

	// Verificar se a empresa existe
	var existingCompany models.Company
	err = database.DB.Get(&existingCompany, "SELECT * FROM companies WHERE id = $1", companyID)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Empresa não encontrada",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao buscar empresa",
		})
	}

	// Construir query de update dinamicamente
	updates := []string{}
	args := []interface{}{}
	argCount := 1

	if req.Name != nil {
		// Verificar se já existe outra empresa com o mesmo nome
		var nameCheckCompany models.Company
		err := database.DB.Get(&nameCheckCompany, "SELECT id FROM companies WHERE name = $1 AND id != $2", *req.Name, companyID)
		if err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"status":  "error",
				"message": "Já existe uma empresa com esse nome",
			})
		} else if err != sql.ErrNoRows {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Erro interno do servidor",
			})
		}

		updates = append(updates, "name = $"+fmt.Sprintf("%d", argCount))
		args = append(args, *req.Name)
		argCount++
	}

	if req.Description != nil {
		updates = append(updates, "description = $"+fmt.Sprintf("%d", argCount))
		args = append(args, *req.Description)
		argCount++
	}

	if req.IsActive != nil {
		updates = append(updates, "is_active = $"+fmt.Sprintf("%d", argCount))
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
	updates = append(updates, "updated_at = $"+fmt.Sprintf("%d", argCount))
	args = append(args, time.Now())
	argCount++

	// Adicionar ID para WHERE clause
	args = append(args, companyID)

	query := fmt.Sprintf("UPDATE companies SET %s WHERE id = $%d", strings.Join(updates, ", "), argCount)
	
	_, err = database.DB.Exec(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao atualizar empresa",
		})
	}

	// Buscar empresa atualizada
	var updatedCompany models.Company
	err = database.DB.Get(&updatedCompany, "SELECT * FROM companies WHERE id = $1", companyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao buscar empresa atualizada",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   updatedCompany,
	})
}

func DeleteCompany(c *fiber.Ctx) error {
	id := c.Params("id")
	companyID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "ID da empresa inválido",
		})
	}

	// Verificar se a empresa existe
	var company models.Company
	err = database.DB.Get(&company, "SELECT id FROM companies WHERE id = $1", companyID)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Empresa não encontrada",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao buscar empresa",
		})
	}

	// Verificar se existem usuários associados à empresa
	var userCount int
	err = database.DB.Get(&userCount, "SELECT COUNT(*) FROM users WHERE company_id = $1", companyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao verificar usuários associados",
		})
	}

	if userCount > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "error",
			"message": "Não é possível excluir uma empresa que possui usuários associados",
		})
	}

	// Deletar empresa
	_, err = database.DB.Exec("DELETE FROM companies WHERE id = $1", companyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Erro ao excluir empresa",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Empresa excluída com sucesso",
	})
}
