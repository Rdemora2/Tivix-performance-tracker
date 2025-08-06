package migrations

// migration001SQL - Configuração inicial PostgreSQL
const migration001SQL = `
-- Habilitar extensão UUID para gerar IDs únicos
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Verificar se a extensão foi criada com sucesso
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp') THEN
        RAISE EXCEPTION 'Falha ao criar extensão uuid-ossp';
    END IF;
    RAISE NOTICE 'Extensão uuid-ossp habilitada com sucesso';
END $$;
`

// migration002SQL - Criação das tabelas principais
const migration002SQL = `
-- Tabela de empresas (multitenant)
CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de usuários
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'manager', 'user')),
    company_id UUID REFERENCES companies(id) ON DELETE SET NULL,
    needs_password_change BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de times/equipes
CREATE TABLE IF NOT EXISTS teams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    color VARCHAR(50) DEFAULT 'blue',
    company_id UUID REFERENCES companies(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de desenvolvedores
CREATE TABLE IF NOT EXISTS developers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL,
    latest_performance_score DECIMAL(4,2) DEFAULT 0.00,
    team_id UUID REFERENCES teams(id) ON DELETE SET NULL,
    company_id UUID REFERENCES companies(id) ON DELETE CASCADE,
    archived_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de relatórios de performance
CREATE TABLE IF NOT EXISTS performance_reports (
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
);
`

// migration003SQL - Criação de índices para performance
const migration003SQL = `
-- Índices para tabela companies
CREATE INDEX IF NOT EXISTS idx_companies_name ON companies(name);
CREATE INDEX IF NOT EXISTS idx_companies_is_active ON companies(is_active);

-- Índices para tabela users
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_users_company_id ON users(company_id);

-- Índices para tabela teams
CREATE INDEX IF NOT EXISTS idx_teams_company_id ON teams(company_id);

-- Índices para tabela developers
CREATE INDEX IF NOT EXISTS idx_developers_team_id ON developers(team_id);
CREATE INDEX IF NOT EXISTS idx_developers_company_id ON developers(company_id);
CREATE INDEX IF NOT EXISTS idx_developers_archived_at ON developers(archived_at);

-- Índices para tabela performance_reports
CREATE INDEX IF NOT EXISTS idx_performance_reports_developer_id ON performance_reports(developer_id);
CREATE INDEX IF NOT EXISTS idx_performance_reports_month ON performance_reports(month);
CREATE INDEX IF NOT EXISTS idx_performance_reports_developer_month ON performance_reports(developer_id, month);
`

// migration004SQL - Configuração de triggers para timestamps
const migration004SQL = `
-- Função para atualizar o campo updated_at automaticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger para tabela companies
DROP TRIGGER IF EXISTS update_companies_updated_at ON companies;
CREATE TRIGGER update_companies_updated_at
    BEFORE UPDATE ON companies
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger para tabela users
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger para tabela teams
DROP TRIGGER IF EXISTS update_teams_updated_at ON teams;
CREATE TRIGGER update_teams_updated_at
    BEFORE UPDATE ON teams
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger para tabela developers
DROP TRIGGER IF EXISTS update_developers_updated_at ON developers;
CREATE TRIGGER update_developers_updated_at
    BEFORE UPDATE ON developers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger para tabela performance_reports
DROP TRIGGER IF EXISTS update_performance_reports_updated_at ON performance_reports;
CREATE TRIGGER update_performance_reports_updated_at
    BEFORE UPDATE ON performance_reports
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
`

// migration005SQL - Implementação do sistema multitenant
const migration005SQL = `
-- Adicionar coluna company_id na tabela users (se não existir)
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='company_id') THEN
        ALTER TABLE users ADD COLUMN company_id UUID REFERENCES companies(id) ON DELETE SET NULL;
    END IF; 
END $$;

-- Adicionar coluna company_id na tabela teams (se não existir)
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='teams' AND column_name='company_id') THEN
        ALTER TABLE teams ADD COLUMN company_id UUID REFERENCES companies(id) ON DELETE CASCADE;
    END IF; 
END $$;

-- Adicionar coluna company_id na tabela developers (se não existir)
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='developers' AND column_name='company_id') THEN
        ALTER TABLE developers ADD COLUMN company_id UUID REFERENCES companies(id) ON DELETE CASCADE;
    END IF; 
END $$;

-- Criar índices adicionais para as novas colunas (se não existirem)
CREATE INDEX IF NOT EXISTS idx_users_company_id ON users(company_id);
CREATE INDEX IF NOT EXISTS idx_teams_company_id ON teams(company_id);
CREATE INDEX IF NOT EXISTS idx_developers_company_id ON developers(company_id);
`

// migration006SQL - Migração de dados para multitenancy
const migration006SQL = `
-- Inserir empresa padrão se não existir
INSERT INTO companies (name, description, is_active, created_at, updated_at) 
SELECT 'Tivix Technologies', 'Empresa padrão do sistema', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM companies WHERE name = 'Tivix Technologies');

-- Atualizar usuários sem empresa
UPDATE users 
SET company_id = (SELECT id FROM companies WHERE name = 'Tivix Technologies' LIMIT 1) 
WHERE company_id IS NULL;

-- Atualizar times sem empresa
UPDATE teams 
SET company_id = (SELECT id FROM companies WHERE name = 'Tivix Technologies' LIMIT 1) 
WHERE company_id IS NULL;

-- Atualizar desenvolvedores sem empresa
UPDATE developers 
SET company_id = (SELECT id FROM companies WHERE name = 'Tivix Technologies' LIMIT 1) 
WHERE company_id IS NULL;

-- Verificar se a migração foi bem-sucedida
DO $$
DECLARE
    users_without_company INTEGER;
    teams_without_company INTEGER;
    developers_without_company INTEGER;
BEGIN
    SELECT COUNT(*) INTO users_without_company FROM users WHERE company_id IS NULL;
    SELECT COUNT(*) INTO teams_without_company FROM teams WHERE company_id IS NULL;
    SELECT COUNT(*) INTO developers_without_company FROM developers WHERE company_id IS NULL;
    
    IF users_without_company > 0 OR teams_without_company > 0 OR developers_without_company > 0 THEN
        RAISE EXCEPTION 'Migração falhou: ainda existem registros sem company_id';
    END IF;
    
    RAISE NOTICE 'Migração de dados para multitenancy concluída com sucesso';
END $$;
`
