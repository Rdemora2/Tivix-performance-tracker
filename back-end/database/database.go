package database

import (
	"fmt"
	"log"

	"tivix-performance-tracker-backend/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Connect() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	var err error
	DB, err = sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("✅ Connected to PostgreSQL database")
}

func Migrate() {
	createTablesQueries := []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,

		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'manager', 'user')),
			needs_password_change BOOLEAN NOT NULL DEFAULT false,
			is_active BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS teams (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			color VARCHAR(50) DEFAULT 'blue',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS developers (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(255) NOT NULL,
			role VARCHAR(255) NOT NULL,
			latest_performance_score DECIMAL(4,2) DEFAULT 0.00,
			team_id UUID REFERENCES teams(id) ON DELETE SET NULL,
			archived_at TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS performance_reports (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			developer_id UUID NOT NULL REFERENCES developers(id) ON DELETE CASCADE,
			month VARCHAR(7) NOT NULL, -- Formato YYYY-MM
			question_scores JSONB NOT NULL,
			category_scores JSONB NOT NULL,
			weighted_average_score DECIMAL(4,2) NOT NULL,
			highlights TEXT,
			points_to_develop TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);`,
		`CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);`,
		`CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);`,
		`CREATE INDEX IF NOT EXISTS idx_developers_team_id ON developers(team_id);`,
		`CREATE INDEX IF NOT EXISTS idx_developers_archived_at ON developers(archived_at);`,
		`CREATE INDEX IF NOT EXISTS idx_performance_reports_developer_id ON performance_reports(developer_id);`,
		`CREATE INDEX IF NOT EXISTS idx_performance_reports_month ON performance_reports(month);`,
		`CREATE INDEX IF NOT EXISTS idx_performance_reports_developer_month ON performance_reports(developer_id, month);`,

		`CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = CURRENT_TIMESTAMP;
			RETURN NEW;
		END;
		$$ language 'plpgsql';`,

		`DROP TRIGGER IF EXISTS update_users_updated_at ON users;
		CREATE TRIGGER update_users_updated_at
			BEFORE UPDATE ON users
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();`,

		`DROP TRIGGER IF EXISTS update_teams_updated_at ON teams;
		CREATE TRIGGER update_teams_updated_at
			BEFORE UPDATE ON teams
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();`,

		`DROP TRIGGER IF EXISTS update_developers_updated_at ON developers;
		CREATE TRIGGER update_developers_updated_at
			BEFORE UPDATE ON developers
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();`,

		`DROP TRIGGER IF EXISTS update_performance_reports_updated_at ON performance_reports;
		CREATE TRIGGER update_performance_reports_updated_at
			BEFORE UPDATE ON performance_reports
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();`,
	}

	for _, query := range createTablesQueries {
		if _, err := DB.Exec(query); err != nil {
			log.Printf("Migration error: %v", err)
			log.Printf("Query: %s", query)
		}
	}

	log.Println("✅ Database migrations completed")
}
