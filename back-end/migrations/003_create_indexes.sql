-- ============================================
-- Migração 003: Criação de Índices para Performance
-- ============================================
-- Descrição: Cria índices para otimizar consultas frequentes
-- Data: 2025-08-05
-- Versão: v1.0.0
-- Autor: Sistema Tivix Performance Tracker
-- ============================================

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
