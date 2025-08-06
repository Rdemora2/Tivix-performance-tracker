-- ============================================
-- Migração 001: Configuração Inicial PostgreSQL
-- ============================================
-- Descrição: Habilita extensões necessárias para o funcionamento do sistema
-- Data: 2025-08-05
-- Versão: v1.0.0
-- Autor: Sistema Tivix Performance Tracker
-- ============================================

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
