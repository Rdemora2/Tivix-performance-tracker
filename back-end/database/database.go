package database

import (
	"fmt"
	"log"

	"tivix-performance-tracker-backend/config"
	"tivix-performance-tracker-backend/migrations"

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
	// Usar o novo sistema de migrações centralizado
	migrationManager := migrations.NewMigrationManager(DB.DB)
	
	log.Println("🔄 Iniciando sistema de migrações centralizado...")
	
	if err := migrationManager.RunMigrations(); err != nil {
		log.Printf("❌ Erro nas migrações: %v", err)
		return
	}
	
	log.Println("✅ Sistema de migrações concluído com sucesso")
}
