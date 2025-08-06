package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"time"
)

type Migration struct {
	ID          string
	Description string
	SQL         string
	AppliedAt   *time.Time
}

type MigrationManager struct {
	DB *sql.DB
}

func NewMigrationManager(db *sql.DB) *MigrationManager {
	return &MigrationManager{DB: db}
}

func (m *MigrationManager) CreateMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id VARCHAR(255) PRIMARY KEY,
			description TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := m.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("falha ao criar tabela de migrações: %w", err)
	}

	log.Println("✅ Tabela de migrações criada/verificada")
	return nil
}

func (m *MigrationManager) GetAppliedMigrations() (map[string]bool, error) {
	applied := make(map[string]bool)

	rows, err := m.DB.Query("SELECT id FROM schema_migrations")
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar migrações aplicadas: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("falha ao ler migração aplicada: %w", err)
		}
		applied[id] = true
	}

	return applied, nil
}

func (m *MigrationManager) RecordMigration(id, description string) error {
	query := `INSERT INTO schema_migrations (id, description) VALUES ($1, $2)`
	_, err := m.DB.Exec(query, id, description)
	if err != nil {
		return fmt.Errorf("falha ao registrar migração %s: %w", id, err)
	}
	return nil
}

func (m *MigrationManager) RunMigrations() error {
	if err := m.CreateMigrationsTable(); err != nil {
		return err
	}

	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	migrations := m.GetAllMigrations()

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].ID < migrations[j].ID
	})

	pendingCount := 0
	for _, migration := range migrations {
		if !applied[migration.ID] {
			log.Printf("🔄 Executando migração %s: %s", migration.ID, migration.Description)

			tx, err := m.DB.Begin()
			if err != nil {
				return fmt.Errorf("falha ao iniciar transação para migração %s: %w", migration.ID, err)
			}

			if _, err := tx.Exec(migration.SQL); err != nil {
				tx.Rollback()
				return fmt.Errorf("falha ao executar migração %s: %w", migration.ID, err)
			}

			if _, err := tx.Exec("INSERT INTO schema_migrations (id, description) VALUES ($1, $2)",
				migration.ID, migration.Description); err != nil {
				tx.Rollback()
				return fmt.Errorf("falha ao registrar migração %s: %w", migration.ID, err)
			}

			if err := tx.Commit(); err != nil {
				return fmt.Errorf("falha ao confirmar migração %s: %w", migration.ID, err)
			}

			log.Printf("✅ Migração %s aplicada com sucesso", migration.ID)
			pendingCount++
		}
	}

	if pendingCount == 0 {
		log.Println("ℹ️  Nenhuma migração pendente encontrada")
	} else {
		log.Printf("✅ %d migração(ões) aplicada(s) com sucesso", pendingCount)
	}

	return nil
}

func (m *MigrationManager) GetAllMigrations() []Migration {
	return []Migration{
		{
			ID:          "001_initial_setup",
			Description: "Configuração inicial PostgreSQL",
			SQL:         migration001SQL,
		},
		{
			ID:          "002_create_tables",
			Description: "Criação das tabelas principais",
			SQL:         migration002SQL,
		},
		{
			ID:          "003_create_indexes",
			Description: "Criação de índices para performance",
			SQL:         migration003SQL,
		},
		{
			ID:          "004_create_triggers",
			Description: "Configuração de triggers para timestamps",
			SQL:         migration004SQL,
		},
		{
			ID:          "005_multitenant_implementation",
			Description: "Implementação do sistema multitenant",
			SQL:         migration005SQL,
		},
		{
			ID:          "006_data_migration_multitenant",
			Description: "Migração de dados para multitenancy",
			SQL:         migration006SQL,
		},
	}
}
