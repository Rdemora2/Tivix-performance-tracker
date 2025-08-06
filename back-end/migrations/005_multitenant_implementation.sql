-- ============================================
-- Migração 005: Implementação do Sistema Multitenant
-- ============================================
-- Descrição: Adiciona colunas company_id para implementar multitenancy
-- Data: 2025-08-05
-- Versão: v1.1.0
-- Autor: Sistema Tivix Performance Tracker
-- ============================================

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
