-- ============================================
-- Migração 006: Migração de Dados para Multitenancy
-- ============================================
-- Descrição: Migra dados existentes para o novo sistema multitenant
-- Data: 2025-08-05
-- Versão: v1.1.0
-- Autor: Sistema Tivix Performance Tracker
-- ============================================

-- Inserir empresa padrão se não existir
INSERT INTO companies (name, description, is_active, created_at, updated_at) 
SELECT 'Tivix Technologies', 'Empresa padrão do sistema', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM companies WHERE name = 'Tivix Technologies');

-- Atualizar registros existentes para vincular à empresa padrão
-- Isso garante que dados existentes não sejam perdidos na migração

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
