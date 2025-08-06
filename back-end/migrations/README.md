# Migrações do Banco de Dados - Tivix Performance Tracker

Este diretório contém todas as migrações do banco de dados para o sistema Tivix Performance Tracker.

## 🔒 Filosofia: Apenas Migrações Automáticas

Este sistema foi projetado para executar migrações **automaticamente** quando a aplicação inicia. Isso garante:

- ✅ **Consistência**: Todas as instâncias sempre têm o mesmo schema
- ✅ **Simplicidade**: Não há comandos manuais para lembrar
- ✅ **Segurança**: Reduz erros humanos em produção
- ✅ **CI/CD Friendly**: Deploys automáticos sem intervenção manual

## Estrutura das Migrações

As migrações são organizadas de forma sequencial e cada arquivo segue a convenção de nomenclatura:

```
{numero}_{descricao}.sql
```

## Histórico de Migrações

| Migração | Descrição                                | Data       | Versão |
| -------- | ---------------------------------------- | ---------- | ------ |
| 001      | Configuração inicial do PostgreSQL       | 2025-08-05 | v1.0.0 |
| 002      | Criação das tabelas principais           | 2025-08-05 | v1.0.0 |
| 003      | Criação de índices para performance      | 2025-08-05 | v1.0.0 |
| 004      | Configuração de triggers para timestamps | 2025-08-05 | v1.0.0 |
| 005      | Implementação do sistema multitenant     | 2025-08-05 | v1.1.0 |
| 006      | Migração de dados para multitenant       | 2025-08-05 | v1.1.0 |

## Como Executar

### Migração Automática

As migrações são executadas **automaticamente** quando a aplicação inicia através do `database.Migrate()`.

```bash
# Iniciar a aplicação (migrações automáticas)
go run main.go
```

### Verificar Status das Migrações

Para verificar quais migrações foram aplicadas:

```bash
# Verificar status
go run cmd/migration-status/main.go
```

### Aplicar Migração Específica (Uso Avançado)

⚠️ **Apenas para desenvolvimento/troubleshooting**

```bash
cd back-end/migrations
psql -h localhost -U seu_usuario -d sua_database -f 001_initial_setup.sql
```

## Regras Importantes

1. **NUNCA modifique migrações já aplicadas em produção**
2. **Sempre crie uma nova migração para mudanças**
3. **Teste migrações em ambiente de desenvolvimento primeiro**
4. **Use transações quando possível**
5. **Documente o propósito de cada migração**

## Schema Atual

### Tabelas Principais

- `companies` - Empresas (sistema multitenant)
- `users` - Usuários do sistema
- `teams` - Times/equipes
- `developers` - Desenvolvedores
- `performance_reports` - Relatórios de performance

### Relacionamentos

- Empresas podem ter múltiplos usuários, times e desenvolvedores
- Times pertencem a uma empresa
- Desenvolvedores pertencem a um time e empresa
- Relatórios de performance são vinculados a desenvolvedores

## Backup e Rollback

Antes de aplicar migrações em produção:

```bash
# Backup completo
pg_dump -h localhost -U usuario -d database > backup_$(date +%Y%m%d_%H%M%S).sql

# Backup apenas do schema
pg_dump -h localhost -U usuario -d database --schema-only > schema_backup_$(date +%Y%m%d_%H%M%S).sql
```

## Troubleshooting

### Erro de Permissões

Certifique-se que o usuário do banco tem permissões para:

- CREATE TABLE
- CREATE INDEX
- CREATE TRIGGER
- CREATE FUNCTION

### Erro de Extensões

O PostgreSQL deve ter a extensão `uuid-ossp` disponível:

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```
