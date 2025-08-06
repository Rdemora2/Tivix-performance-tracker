package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"tivix-performance-tracker-backend/config"
	"tivix-performance-tracker-backend/database"
	"tivix-performance-tracker-backend/migrations"
)

func main() {
	config.LoadConfig()

	database.Connect()

	log.Println("📊 Verificando status das migrações...")

	migrationManager := migrations.NewMigrationManager(database.DB.DB)

	if err := migrationManager.CreateMigrationsTable(); err != nil {
		log.Printf("❌ Erro ao criar tabela de migrações: %v", err)
		return
	}

	applied, err := migrationManager.GetAppliedMigrations()
	if err != nil {
		log.Printf("❌ Erro ao consultar migrações aplicadas: %v", err)
		return
	}

	allMigrations := migrationManager.GetAllMigrations()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tDescrição\tStatus\tData")
	fmt.Fprintln(w, "---\t----------\t------\t----")

	pendingCount := 0
	appliedCount := 0

	for _, migration := range allMigrations {
		status := "⏳ Pendente"
		date := "-"

		if applied[migration.ID] {
			status = "✅ Aplicada"
			appliedCount++

			var appliedAt time.Time
			err := database.DB.QueryRow("SELECT applied_at FROM schema_migrations WHERE id = $1", migration.ID).Scan(&appliedAt)
			if err == nil {
				date = appliedAt.Format("2006-01-02 15:04:05")
			}
		} else {
			pendingCount++
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			migration.ID,
			migration.Description,
			status,
			date)
	}

	w.Flush()

	fmt.Println()
	fmt.Printf("📈 Resumo das Migrações:\n")
	fmt.Printf("   • Total: %d\n", len(allMigrations))
	fmt.Printf("   • Aplicadas: %d\n", appliedCount)
	fmt.Printf("   • Pendentes: %d\n", pendingCount)

	if pendingCount > 0 {
		fmt.Println()
		fmt.Println("⚠️  Existem migrações pendentes.")
		fmt.Println("   Inicie a aplicação para aplicá-las automaticamente:")
		fmt.Println("   go run main.go")
	} else {
		fmt.Println()
		fmt.Println("🎉 Todas as migrações estão atualizadas!")
	}
}
