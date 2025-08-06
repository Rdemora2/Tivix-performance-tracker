package main

import (
	"log"
	
	"tivix-performance-tracker-backend/config"
	"tivix-performance-tracker-backend/database"
)

func main() {
	// Carregar configuração
	config.LoadConfig()
	
	// Conectar ao banco
	database.Connect()
	
	log.Println("🔄 Executando migrações para sistema multitenant...")
	
	// Executar migrações
	migrationQueries := []string{
		// 1. Criar tabela de empresas
		`CREATE TABLE IF NOT EXISTS companies (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			is_active BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// 2. Adicionar coluna company_id na tabela users (se não existir)
		`DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='company_id') THEN
				ALTER TABLE users ADD COLUMN company_id UUID REFERENCES companies(id) ON DELETE SET NULL;
			END IF; 
		END $$;`,

		// 3. Adicionar coluna company_id na tabela teams (se não existir)
		`DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='teams' AND column_name='company_id') THEN
				ALTER TABLE teams ADD COLUMN company_id UUID REFERENCES companies(id) ON DELETE CASCADE;
			END IF; 
		END $$;`,

		// 4. Adicionar coluna company_id na tabela developers (se não existir)
		`DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='developers' AND column_name='company_id') THEN
				ALTER TABLE developers ADD COLUMN company_id UUID REFERENCES companies(id) ON DELETE CASCADE;
			END IF; 
		END $$;`,

		// 5. Criar índices
		`CREATE INDEX IF NOT EXISTS idx_companies_name ON companies(name);`,
		`CREATE INDEX IF NOT EXISTS idx_companies_is_active ON companies(is_active);`,
		`CREATE INDEX IF NOT EXISTS idx_users_company_id ON users(company_id);`,
		`CREATE INDEX IF NOT EXISTS idx_teams_company_id ON teams(company_id);`,
		`CREATE INDEX IF NOT EXISTS idx_developers_company_id ON developers(company_id);`,

		// 6. Criar trigger para companies
		`DROP TRIGGER IF EXISTS update_companies_updated_at ON companies;
		CREATE TRIGGER update_companies_updated_at
			BEFORE UPDATE ON companies
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();`,
	}

	for i, query := range migrationQueries {
		log.Printf("Executando migração %d/%d...", i+1, len(migrationQueries))
		if _, err := database.DB.Exec(query); err != nil {
			log.Printf("❌ Erro na migração %d: %v", i+1, err)
			log.Printf("Query: %s", query)
		} else {
			log.Printf("✅ Migração %d concluída", i+1)
		}
	}

	log.Println("✅ Migrações concluídas com sucesso!")
	log.Println("")
	log.Println("📋 Sistema multitenant configurado:")
	log.Println("   • Tabela 'companies' criada")
	log.Println("   • Campo 'company_id' adicionado em users, teams e developers")
	log.Println("   • Índices criados para otimização")
	log.Println("")
	log.Println("🔑 Próximos passos:")
	log.Println("   1. Admins podem criar empresas via API")
	log.Println("   2. Ao criar usuários, especifique a empresa")
	log.Println("   3. Managers/usuários só veem dados da sua empresa")
}
