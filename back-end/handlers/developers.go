package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"tivix-performance-tracker-backend/database"
	"tivix-performance-tracker-backend/middleware"
	"tivix-performance-tracker-backend/models"
)

// GetAllDevelopers retorna todos os desenvolvedores
func GetAllDevelopers(c *fiber.Ctx) error {
	user := c.Locals("user").(*middleware.JWTClaims)
	
	// Verificar se deve incluir arquivados
	includeArchived := c.Query("includeArchived", "false")

	var query string
	var args []interface{}

	if user.Role == "admin" {
		// Admins podem ver todos os desenvolvedores
		query = `
			SELECT id, name, role, latest_performance_score, team_id, company_id, archived_at, created_at, updated_at 
			FROM developers 
		`
		if includeArchived != "true" {
			query += " WHERE archived_at IS NULL"
		}
	} else {
		// Managers e usuários só podem ver desenvolvedores da sua empresa
		if user.CompanyID == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": "Usuário deve estar associado a uma empresa",
			})
		}
		
		query = `
			SELECT id, name, role, latest_performance_score, team_id, company_id, archived_at, created_at, updated_at 
			FROM developers 
			WHERE company_id = $1
		`
		args = append(args, *user.CompanyID)
		
		if includeArchived != "true" {
			query += " AND archived_at IS NULL"
		}
	}

	query += " ORDER BY created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error querying developers: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao buscar desenvolvedores",
		})
	}
	defer rows.Close()

	var developers []models.Developer
	for rows.Next() {
		var developer models.Developer
		err := rows.Scan(
			&developer.ID,
			&developer.Name,
			&developer.Role,
			&developer.LatestPerformanceScore,
			&developer.TeamID,
			&developer.CompanyID,
			&developer.ArchivedAt,
			&developer.CreatedAt,
			&developer.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning developer: %v", err)
			continue
		}
		developers = append(developers, developer)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    developers,
	})
}

// GetArchivedDevelopers retorna apenas desenvolvedores arquivados
func GetArchivedDevelopers(c *fiber.Ctx) error {
	query := `
		SELECT id, name, role, latest_performance_score, team_id, archived_at, created_at, updated_at 
		FROM developers 
		WHERE archived_at IS NOT NULL
		ORDER BY archived_at DESC
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Error querying archived developers: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao buscar desenvolvedores arquivados",
		})
	}
	defer rows.Close()

	var developers []models.Developer
	for rows.Next() {
		var developer models.Developer
		err := rows.Scan(
			&developer.ID,
			&developer.Name,
			&developer.Role,
			&developer.LatestPerformanceScore,
			&developer.TeamID,
			&developer.ArchivedAt,
			&developer.CreatedAt,
			&developer.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning developer: %v", err)
			continue
		}
		developers = append(developers, developer)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    developers,
	})
}

// GetDeveloperByID retorna um desenvolvedor específico por ID
func GetDeveloperByID(c *fiber.Ctx) error {
	id := c.Params("id")
	developerUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "ID inválido",
		})
	}

	query := `
		SELECT id, name, role, latest_performance_score, team_id, archived_at, created_at, updated_at 
		FROM developers 
		WHERE id = $1
	`

	var developer models.Developer
	err = database.DB.QueryRow(query, developerUUID).Scan(
		&developer.ID,
		&developer.Name,
		&developer.Role,
		&developer.LatestPerformanceScore,
		&developer.TeamID,
		&developer.ArchivedAt,
		&developer.CreatedAt,
		&developer.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "Desenvolvedor não encontrado",
		})
	}
	if err != nil {
		log.Printf("Error querying developer: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao buscar desenvolvedor",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    developer,
	})
}

// CreateDeveloper cria um novo desenvolvedor
func CreateDeveloper(c *fiber.Ctx) error {
	user := c.Locals("user").(*middleware.JWTClaims)
	
	var req models.CreateDeveloperRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Dados inválidos",
		})
	}

	if req.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Nome é obrigatório",
		})
	}

	if req.Role == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Cargo é obrigatório",
		})
	}

	// Determinar a empresa do desenvolvedor
	var companyID *uuid.UUID
	if user.Role == "admin" && req.CompanyID != nil {
		// Admin pode especificar a empresa
		companyID = req.CompanyID
	} else if user.CompanyID != nil {
		// Managers e usuários criam desenvolvedores na sua própria empresa
		companyID = user.CompanyID
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   true,
			"message": "Usuário deve estar associado a uma empresa",
		})
	}

	// Verificar se o team_id existe e pertence à mesma empresa (se fornecido)
	if req.TeamID != nil {
		var teamCompanyID *uuid.UUID
		err := database.DB.QueryRow("SELECT company_id FROM teams WHERE id = $1", *req.TeamID).Scan(&teamCompanyID)
		if err == sql.ErrNoRows {
			return c.Status(400).JSON(fiber.Map{
				"error":   true,
				"message": "Time não encontrado",
			})
		} else if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   true,
				"message": "Erro ao verificar time",
			})
		}

		// Verificar se o time pertence à mesma empresa
		if user.Role != "admin" && (teamCompanyID == nil || companyID == nil || *teamCompanyID != *companyID) {
			return c.Status(400).JSON(fiber.Map{
				"error":   true,
				"message": "Time não pertence à sua empresa",
			})
		}
	}

	query := `
		INSERT INTO developers (name, role, team_id, company_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, role, latest_performance_score, team_id, company_id, archived_at, created_at, updated_at
	`

	var developer models.Developer
	err := database.DB.QueryRow(query, req.Name, req.Role, req.TeamID, companyID).Scan(
		&developer.ID,
		&developer.Name,
		&developer.Role,
		&developer.LatestPerformanceScore,
		&developer.TeamID,
		&developer.CompanyID,
		&developer.ArchivedAt,
		&developer.CreatedAt,
		&developer.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error creating developer: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao criar desenvolvedor",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    developer,
	})
}

// UpdateDeveloper atualiza um desenvolvedor existente
func UpdateDeveloper(c *fiber.Ctx) error {
	id := c.Params("id")
	developerUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "ID inválido",
		})
	}

	var req models.UpdateDeveloperRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Dados inválidos",
		})
	}

	// Verificar se o desenvolvedor existe
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM developers WHERE id = $1)", developerUUID).Scan(&exists)
	if err != nil || !exists {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "Desenvolvedor não encontrado",
		})
	}

	// Verificar se o team_id existe (se fornecido)
	if req.TeamID != nil {
		var teamExists bool
		err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM teams WHERE id = $1)", *req.TeamID).Scan(&teamExists)
		if err != nil || !teamExists {
			return c.Status(400).JSON(fiber.Map{
				"error":   true,
				"message": "Time não encontrado",
			})
		}
	}

	// Construir query dinâmica
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Role != nil {
		setParts = append(setParts, fmt.Sprintf("role = $%d", argIndex))
		args = append(args, *req.Role)
		argIndex++
	}
	if req.LatestPerformanceScore != nil {
		setParts = append(setParts, fmt.Sprintf("latest_performance_score = $%d", argIndex))
		args = append(args, *req.LatestPerformanceScore)
		argIndex++
	}
	if req.TeamID != nil {
		setParts = append(setParts, fmt.Sprintf("team_id = $%d", argIndex))
		args = append(args, *req.TeamID)
		argIndex++
	}

	if len(setParts) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Nenhum campo para atualizar",
		})
	}

	query := "UPDATE developers SET "
	for i, part := range setParts {
		if i > 0 {
			query += ", "
		}
		query += part
	}
	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, name, role, latest_performance_score, team_id, archived_at, created_at, updated_at", argIndex)

	args = append(args, developerUUID)

	var developer models.Developer
	err = database.DB.QueryRow(query, args...).Scan(
		&developer.ID,
		&developer.Name,
		&developer.Role,
		&developer.LatestPerformanceScore,
		&developer.TeamID,
		&developer.ArchivedAt,
		&developer.CreatedAt,
		&developer.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error updating developer: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao atualizar desenvolvedor",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    developer,
	})
}

// ArchiveDeveloper arquiva ou restaura um desenvolvedor
func ArchiveDeveloper(c *fiber.Ctx) error {
	id := c.Params("id")
	developerUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "ID inválido",
		})
	}

	var req models.ArchiveDeveloperRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Dados inválidos",
		})
	}

	var query string
	var archivedAt *time.Time

	if req.Archive {
		// Arquivar desenvolvedor
		now := time.Now()
		archivedAt = &now
		query = `
			UPDATE developers 
			SET archived_at = $1 
			WHERE id = $2 
			RETURNING id, name, role, latest_performance_score, team_id, archived_at, created_at, updated_at
		`
	} else {
		// Restaurar desenvolvedor
		query = `
			UPDATE developers 
			SET archived_at = NULL 
			WHERE id = $1 
			RETURNING id, name, role, latest_performance_score, team_id, archived_at, created_at, updated_at
		`
	}

	var developer models.Developer
	var scanErr error

	if req.Archive {
		scanErr = database.DB.QueryRow(query, archivedAt, developerUUID).Scan(
			&developer.ID,
			&developer.Name,
			&developer.Role,
			&developer.LatestPerformanceScore,
			&developer.TeamID,
			&developer.ArchivedAt,
			&developer.CreatedAt,
			&developer.UpdatedAt,
		)
	} else {
		scanErr = database.DB.QueryRow(query, developerUUID).Scan(
			&developer.ID,
			&developer.Name,
			&developer.Role,
			&developer.LatestPerformanceScore,
			&developer.TeamID,
			&developer.ArchivedAt,
			&developer.CreatedAt,
			&developer.UpdatedAt,
		)
	}

	if scanErr == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "Desenvolvedor não encontrado",
		})
	}
	if scanErr != nil {
		log.Printf("Error archiving/restoring developer: %v", scanErr)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao arquivar/restaurar desenvolvedor",
		})
	}

	action := "restaurado"
	if req.Archive {
		action = "arquivado"
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Desenvolvedor %s com sucesso", action),
		"data":    developer,
	})
}

// GetDevelopersByTeam retorna desenvolvedores de um time específico
func GetDevelopersByTeam(c *fiber.Ctx) error {
	teamID := c.Params("teamId")
	teamUUID, err := uuid.Parse(teamID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "ID do time inválido",
		})
	}

	// Verificar se deve incluir arquivados
	includeArchived := c.Query("includeArchived", "false")

	query := `
		SELECT id, name, role, latest_performance_score, team_id, archived_at, created_at, updated_at 
		FROM developers 
		WHERE team_id = $1
	`

	if includeArchived != "true" {
		query += " AND archived_at IS NULL"
	}

	query += " ORDER BY created_at DESC"

	rows, err := database.DB.Query(query, teamUUID)
	if err != nil {
		log.Printf("Error querying developers by team: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao buscar desenvolvedores do time",
		})
	}
	defer rows.Close()

	var developers []models.Developer
	for rows.Next() {
		var developer models.Developer
		err := rows.Scan(
			&developer.ID,
			&developer.Name,
			&developer.Role,
			&developer.LatestPerformanceScore,
			&developer.TeamID,
			&developer.ArchivedAt,
			&developer.CreatedAt,
			&developer.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning developer: %v", err)
			continue
		}
		developers = append(developers, developer)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    developers,
	})
}

func DeleteDeveloper(c *fiber.Ctx) error {
	id := c.Params("id")
	developerUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "ID inválido",
		})
	}

	// Obter usuário atual das claims do JWT
	user := c.Locals("user").(*middleware.JWTClaims)

	var existingDeveloper models.Developer
	checkQuery := `
		SELECT id, name, role, latest_performance_score, team_id, company_id, archived_at, created_at, updated_at 
		FROM developers 
		WHERE id = $1
	`

	err = database.DB.QueryRow(checkQuery, developerUUID).Scan(
		&existingDeveloper.ID,
		&existingDeveloper.Name,
		&existingDeveloper.Role,
		&existingDeveloper.LatestPerformanceScore,
		&existingDeveloper.TeamID,
		&existingDeveloper.CompanyID,
		&existingDeveloper.ArchivedAt,
		&existingDeveloper.CreatedAt,
		&existingDeveloper.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "Desenvolvedor não encontrado",
		})
	}
	if err != nil {
		log.Printf("Error checking developer existence: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao verificar desenvolvedor",
		})
	}

	// Verificar permissões de exclusão para managers
	if user.Role != "admin" {
		// Managers só podem excluir desenvolvedores da sua própria empresa
		if user.CompanyID == nil || existingDeveloper.CompanyID == nil || *user.CompanyID != *existingDeveloper.CompanyID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": "Sem permissão para excluir este desenvolvedor",
			})
		}
	}

	// Inicia uma transação para garantir consistência
	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro interno do servidor",
		})
	}
	defer tx.Rollback()

	// Primeiro, exclui todos os relatórios de performance do desenvolvedor
	_, err = tx.Exec("DELETE FROM performance_reports WHERE developer_id = $1", developerUUID)
	if err != nil {
		log.Printf("Error deleting performance reports: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao excluir relatórios de performance",
		})
	}

	// Agora exclui o desenvolvedor
	result, err := tx.Exec("DELETE FROM developers WHERE id = $1", developerUUID)
	if err != nil {
		log.Printf("Error deleting developer: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao excluir desenvolvedor",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "Desenvolvedor não encontrado",
		})
	}

	// Confirma a transação
	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao confirmar exclusão",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Desenvolvedor excluído com sucesso",
		"data": fiber.Map{
			"deletedDeveloper": existingDeveloper,
		},
	})
}
