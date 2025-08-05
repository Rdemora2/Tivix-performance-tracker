package handlers

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"tivix-performance-tracker-backend/database"
	"tivix-performance-tracker-backend/models"
)

// GetAllPerformanceReports retorna todos os relatórios de performance
func GetAllPerformanceReports(c *fiber.Ctx) error {
	query := `
		SELECT pr.id, pr.developer_id, pr.month, pr.question_scores, pr.category_scores, 
		       pr.weighted_average_score, pr.highlights, pr.points_to_develop, 
		       pr.created_at, pr.updated_at
		FROM performance_reports pr
		ORDER BY pr.month DESC, pr.created_at DESC
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Error querying performance reports: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao buscar relatórios de performance",
		})
	}
	defer rows.Close()

	var reports []models.PerformanceReport
	for rows.Next() {
		var report models.PerformanceReport
		err := rows.Scan(
			&report.ID,
			&report.DeveloperID,
			&report.Month,
			&report.QuestionScores,
			&report.CategoryScores,
			&report.WeightedAverageScore,
			&report.Highlights,
			&report.PointsToDevelop,
			&report.CreatedAt,
			&report.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning performance report: %v", err)
			continue
		}
		reports = append(reports, report)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    reports,
	})
}

// GetPerformanceReportsByDeveloper retorna relatórios de performance de um desenvolvedor
func GetPerformanceReportsByDeveloper(c *fiber.Ctx) error {
	developerID := c.Params("developerId")
	developerUUID, err := uuid.Parse(developerID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "ID do desenvolvedor inválido",
		})
	}

	query := `
		SELECT id, developer_id, month, question_scores, category_scores, 
		       weighted_average_score, highlights, points_to_develop, 
		       created_at, updated_at
		FROM performance_reports 
		WHERE developer_id = $1
		ORDER BY month DESC, created_at DESC
	`

	rows, err := database.DB.Query(query, developerUUID)
	if err != nil {
		log.Printf("Error querying performance reports by developer: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao buscar relatórios do desenvolvedor",
		})
	}
	defer rows.Close()

	var reports []models.PerformanceReport
	for rows.Next() {
		var report models.PerformanceReport
		err := rows.Scan(
			&report.ID,
			&report.DeveloperID,
			&report.Month,
			&report.QuestionScores,
			&report.CategoryScores,
			&report.WeightedAverageScore,
			&report.Highlights,
			&report.PointsToDevelop,
			&report.CreatedAt,
			&report.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning performance report: %v", err)
			continue
		}
		reports = append(reports, report)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    reports,
	})
}

// GetPerformanceReportsByMonth retorna relatórios de performance de um mês específico
func GetPerformanceReportsByMonth(c *fiber.Ctx) error {
	month := c.Params("month")

	query := `
		SELECT pr.id, pr.developer_id, pr.month, pr.question_scores, pr.category_scores, 
		       pr.weighted_average_score, pr.highlights, pr.points_to_develop, 
		       pr.created_at, pr.updated_at
		FROM performance_reports pr
		WHERE pr.month = $1
		ORDER BY pr.weighted_average_score DESC, pr.created_at DESC
	`

	rows, err := database.DB.Query(query, month)
	if err != nil {
		log.Printf("Error querying performance reports by month: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao buscar relatórios do mês",
		})
	}
	defer rows.Close()

	var reports []models.PerformanceReport
	for rows.Next() {
		var report models.PerformanceReport
		err := rows.Scan(
			&report.ID,
			&report.DeveloperID,
			&report.Month,
			&report.QuestionScores,
			&report.CategoryScores,
			&report.WeightedAverageScore,
			&report.Highlights,
			&report.PointsToDevelop,
			&report.CreatedAt,
			&report.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning performance report: %v", err)
			continue
		}
		reports = append(reports, report)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    reports,
	})
}

// GetPerformanceReportByID retorna um relatório específico por ID
func GetPerformanceReportByID(c *fiber.Ctx) error {
	id := c.Params("id")
	reportUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "ID inválido",
		})
	}

	query := `
		SELECT id, developer_id, month, question_scores, category_scores, 
		       weighted_average_score, highlights, points_to_develop, 
		       created_at, updated_at
		FROM performance_reports 
		WHERE id = $1
	`

	var report models.PerformanceReport
	err = database.DB.QueryRow(query, reportUUID).Scan(
		&report.ID,
		&report.DeveloperID,
		&report.Month,
		&report.QuestionScores,
		&report.CategoryScores,
		&report.WeightedAverageScore,
		&report.Highlights,
		&report.PointsToDevelop,
		&report.CreatedAt,
		&report.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "Relatório não encontrado",
		})
	}
	if err != nil {
		log.Printf("Error querying performance report: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao buscar relatório",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    report,
	})
}

// CreatePerformanceReport cria um novo relatório de performance
func CreatePerformanceReport(c *fiber.Ctx) error {
	var req models.CreatePerformanceReportRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Dados inválidos",
		})
	}

	// Validações básicas
	if req.DeveloperID == uuid.Nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "ID do desenvolvedor é obrigatório",
		})
	}

	if req.Month == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Mês é obrigatório",
		})
	}

	if req.WeightedAverageScore < 0 || req.WeightedAverageScore > 10 {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Pontuação deve estar entre 0 e 10",
		})
	}

	// Verificar se o desenvolvedor existe
	var developerExists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM developers WHERE id = $1)", req.DeveloperID).Scan(&developerExists)
	if err != nil || !developerExists {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Desenvolvedor não encontrado",
		})
	}

	// Verificar se já existe um relatório para este desenvolvedor neste mês
	var existingReportExists bool
	err = database.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM performance_reports WHERE developer_id = $1 AND month = $2)",
		req.DeveloperID, req.Month,
	).Scan(&existingReportExists)
	if err != nil {
		log.Printf("Error checking existing report: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao verificar relatório existente",
		})
	}
	if existingReportExists {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Já existe um relatório para este desenvolvedor neste mês",
		})
	}

	// Inserir novo relatório
	query := `
		INSERT INTO performance_reports (developer_id, month, question_scores, category_scores, 
		                               weighted_average_score, highlights, points_to_develop)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, developer_id, month, question_scores, category_scores, 
		          weighted_average_score, highlights, points_to_develop, created_at, updated_at
	`

	var report models.PerformanceReport
	err = database.DB.QueryRow(
		query,
		req.DeveloperID,
		req.Month,
		req.QuestionScores,
		req.CategoryScores,
		req.WeightedAverageScore,
		req.Highlights,
		req.PointsToDevelop,
	).Scan(
		&report.ID,
		&report.DeveloperID,
		&report.Month,
		&report.QuestionScores,
		&report.CategoryScores,
		&report.WeightedAverageScore,
		&report.Highlights,
		&report.PointsToDevelop,
		&report.CreatedAt,
		&report.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error creating performance report: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao criar relatório",
		})
	}

	// Atualizar a pontuação mais recente do desenvolvedor
	_, err = database.DB.Exec(
		"UPDATE developers SET latest_performance_score = $1 WHERE id = $2",
		req.WeightedAverageScore,
		req.DeveloperID,
	)
	if err != nil {
		log.Printf("Error updating developer latest score: %v", err)
		// Não retorna erro porque o relatório foi criado com sucesso
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    report,
	})
}

// GetAvailableMonths retorna os meses disponíveis com relatórios
func GetAvailableMonths(c *fiber.Ctx) error {
	query := `
		SELECT DISTINCT month 
		FROM performance_reports 
		ORDER BY month DESC
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Error querying available months: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao buscar meses disponíveis",
		})
	}
	defer rows.Close()

	var months []string
	for rows.Next() {
		var month string
		err := rows.Scan(&month)
		if err != nil {
			log.Printf("Error scanning month: %v", err)
			continue
		}
		months = append(months, month)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    months,
	})
}

// GetPerformanceStats retorna estatísticas gerais de performance
func GetPerformanceStats(c *fiber.Ctx) error {
	query := `
		SELECT 
			COUNT(*) as total_reports,
			ROUND(AVG(weighted_average_score)::numeric, 2) as average_score,
			MAX(weighted_average_score) as highest_score,
			MIN(weighted_average_score) as lowest_score
		FROM performance_reports
	`

	var stats struct {
		TotalReports  int     `json:"totalReports"`
		AverageScore  float64 `json:"averageScore"`
		HighestScore  float64 `json:"highestScore"`
		LowestScore   float64 `json:"lowestScore"`
	}

	err := database.DB.QueryRow(query).Scan(
		&stats.TotalReports,
		&stats.AverageScore,
		&stats.HighestScore,
		&stats.LowestScore,
	)

	if err != nil {
		log.Printf("Error querying performance stats: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Erro ao buscar estatísticas",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}
