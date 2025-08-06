-- ============================================
-- Migração 004: Configuração de Triggers para Timestamps
-- ============================================
-- Descrição: Cria função e triggers para atualização automática de updated_at
-- Data: 2025-08-05
-- Versão: v1.0.0
-- Autor: Sistema Tivix Performance Tracker
-- ============================================

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
