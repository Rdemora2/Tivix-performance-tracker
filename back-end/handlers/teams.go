package handlers

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"tivix-performance-tracker-backend/database"
	"tivix-performance-tracker-backend/models"
)

// GetAllTeams retorna todos os times
func GetAllTeams(c *fiber.Ctx) error {
	query := `
		SELECT id, name, description, color, created_at, updated_at 
		FROM teams 
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Error querying teams: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao buscar times",
		})
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var team models.Team
		err := rows.Scan(
			&team.ID,
			&team.Name,
			&team.Description,
			&team.Color,
			&team.CreatedAt,
			&team.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning team: %v", err)
			continue
		}
		teams = append(teams, team)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    teams,
	})
}

// GetTeamByID retorna um time específico por ID
func GetTeamByID(c *fiber.Ctx) error {
	id := c.Params("id")
	teamUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "ID inválido",
		})
	}

	query := `
		SELECT id, name, description, color, created_at, updated_at 
		FROM teams 
		WHERE id = $1
	`

	var team models.Team
	err = database.DB.QueryRow(query, teamUUID).Scan(
		&team.ID,
		&team.Name,
		&team.Description,
		&team.Color,
		&team.CreatedAt,
		&team.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "Time não encontrado",
		})
	}
	if err != nil {
		log.Printf("Error querying team: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao buscar time",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    team,
	})
}

// CreateTeam cria um novo time
func CreateTeam(c *fiber.Ctx) error {
	var req models.CreateTeamRequest
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

	if req.Color == "" {
		req.Color = "blue"
	}

	query := `
		INSERT INTO teams (name, description, color)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, color, created_at, updated_at
	`

	var team models.Team
	err := database.DB.QueryRow(query, req.Name, req.Description, req.Color).Scan(
		&team.ID,
		&team.Name,
		&team.Description,
		&team.Color,
		&team.CreatedAt,
		&team.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error creating team: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao criar time",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    team,
	})
}

// UpdateTeam atualiza um time existente
func UpdateTeam(c *fiber.Ctx) error {
	id := c.Params("id")
	teamUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "ID inválido",
		})
	}

	var req models.UpdateTeamRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Dados inválidos",
		})
	}

	// Verificar se o time existe
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM teams WHERE id = $1)", teamUUID).Scan(&exists)
	if err != nil || !exists {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "Time não encontrado",
		})
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
	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}
	if req.Color != nil {
		setParts = append(setParts, fmt.Sprintf("color = $%d", argIndex))
		args = append(args, *req.Color)
		argIndex++
	}

	if len(setParts) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Nenhum campo para atualizar",
		})
	}

	query := "UPDATE teams SET "
	for i, part := range setParts {
		if i > 0 {
			query += ", "
		}
		query += part
	}
	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, name, description, color, created_at, updated_at", argIndex)

	args = append(args, teamUUID)

	var team models.Team
	err = database.DB.QueryRow(query, args...).Scan(
		&team.ID,
		&team.Name,
		&team.Description,
		&team.Color,
		&team.CreatedAt,
		&team.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error updating team: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao atualizar time",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    team,
	})
}

// DeleteTeam exclui um time
func DeleteTeam(c *fiber.Ctx) error {
	id := c.Params("id")
	teamUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "ID inválido",
		})
	}

	// Primeiro, remove a associação dos desenvolvedores com o time
	_, err = database.DB.Exec("UPDATE developers SET team_id = NULL WHERE team_id = $1", teamUUID)
	if err != nil {
		log.Printf("Error removing team associations: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao remover associações do time",
		})
	}

	// Agora exclui o time
	result, err := database.DB.Exec("DELETE FROM teams WHERE id = $1", teamUUID)
	if err != nil {
		log.Printf("Error deleting team: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao excluir time",
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "Time não encontrado",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Time excluído com sucesso",
	})
}
