# Tivix Performance Tracker Backend

Backend em Go com Fiber para o sistema de avaliação de performance da Tivix.

## 🚀 Tecnologias

- **Go 1.21+** - Linguagem de programação
- **Fiber v2** - Framework web rápido e expressivo
- **PostgreSQL** - Banco de dados relacional
- **UUID** - Identificadores únicos
- **CORS** - Cross-Origin Resource Sharing
- **JWT** - Autenticação baseada em tokens
- **bcrypt** - Hash seguro de senhas
- **SQLx** - Driver SQL extendido para Go

## 🔐 Sistema de Autenticação

O sistema inclui autenticação completa baseada em JWT com três níveis de permissão:

- **Admin**: Acesso completo (CRUD em todos os recursos)
- **Manager**: Pode criar/editar desenvolvedores, times e relatórios
- **User**: Apenas visualização

### Inicialização do Sistema

1. **Verificar se o sistema foi inicializado**: `GET /api/v1/init/check`
2. **Criar primeiro usuário admin**: `POST /api/v1/init/admin` (requer chave de instalação)

### Endpoints de Autenticação

- `POST /api/v1/auth/register` - Cadastro de usuário
- `POST /api/v1/auth/login` - Login
- `GET /api/v1/auth/profile` - Perfil do usuário (protegido)
- `POST /api/v1/auth/refresh` - Renovar token (protegido)

📖 **Documentação completa de integração**: Ver `FRONTEND_INTEGRATION.md`

## 📋 Pré-requisitos

- Go 1.21 ou superior
- PostgreSQL
- Arquivo `.env` configurado (localizado na raiz do projeto)

## 🛠️ Instalação

1. **Instalar dependências**:
```bash
cd backend
go mod download
```

2. **Configurar banco de dados**:
   - Certifique-se de que o PostgreSQL está rodando
   - Configure as variáveis de ambiente no arquivo `.env` na raiz do projeto
   - As migrações são executadas automaticamente na inicialização

3. **Executar o servidor**:
```bash
go run main.go
```

O servidor será iniciado na porta `8080` por padrão.

## 🌐 Endpoints da API

### Times (`/api/v1/teams`)

- `GET /` - Listar todos os times
- `GET /:id` - Buscar time por ID
- `POST /` - Criar novo time
- `PUT /:id` - Atualizar time
- `DELETE /:id` - Excluir time
- `GET /:teamId/developers` - Listar desenvolvedores do time

### Desenvolvedores (`/api/v1/developers`)

- `GET /` - Listar todos os desenvolvedores ativos
- `GET /?includeArchived=true` - Listar incluindo arquivados
- `GET /archived` - Listar apenas arquivados
- `GET /:id` - Buscar desenvolvedor por ID
- `POST /` - Criar novo desenvolvedor
- `PUT /:id` - Atualizar desenvolvedor
- `PUT /:id/archive` - Arquivar/restaurar desenvolvedor
- `GET /:developerId/reports` - Listar relatórios do desenvolvedor

### Relatórios de Performance (`/api/v1/performance-reports`)

- `GET /` - Listar todos os relatórios
- `GET /:id` - Buscar relatório por ID
- `POST /` - Criar novo relatório
- `GET /months` - Listar meses disponíveis
- `GET /month/:month` - Relatórios de um mês específico
- `GET /stats` - Estatísticas gerais

## 📊 Estrutura do Banco de Dados

### Tabela `teams`
- `id` (UUID, PK)
- `name` (VARCHAR)
- `description` (TEXT)
- `color` (VARCHAR)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### Tabela `developers`
- `id` (UUID, PK)
- `name` (VARCHAR)
- `role` (VARCHAR)
- `latest_performance_score` (DECIMAL)
- `team_id` (UUID, FK)
- `archived_at` (TIMESTAMP, nullable)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### Tabela `performance_reports`
- `id` (UUID, PK)
- `developer_id` (UUID, FK)
- `month` (VARCHAR) - Formato YYYY-MM
- `question_scores` (JSONB)
- `category_scores` (JSONB)
- `weighted_average_score` (DECIMAL)
- `highlights` (TEXT)
- `points_to_develop` (TEXT)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

## 🔧 Variáveis de Ambiente

```env
# Database Configuration
VITE_DB_HOST=localhost
VITE_DB_PORT=5433
VITE_DB_USER=perf_tracker
VITE_DB_PASSWORD=your_password
VITE_DB_NAME=performance_tracker_db

# Server Configuration
PORT=8080
```

## 📝 Exemplo de Uso

### Criar um time
```bash
curl -X POST http://localhost:8080/api/v1/teams \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Frontend Team",
    "description": "Equipe responsável pelo desenvolvimento frontend",
    "color": "blue"
  }'
```

### Criar um desenvolvedor
```bash
curl -X POST http://localhost:8080/api/v1/developers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "role": "Frontend Developer",
    "teamId": "uuid-do-time"
  }'
```

### Criar um relatório de performance
```bash
curl -X POST http://localhost:8080/api/v1/performance-reports \
  -H "Content-Type: application/json" \
  -d '{
    "developerId": "uuid-do-desenvolvedor",
    "month": "2024-12",
    "questionScores": {
      "punctualityDeliveries": 8.5,
      "deliveryQuality": 9.0
    },
    "categoryScores": {
      "commitment": 8.5,
      "technicalQuality": 9.0
    },
    "weightedAverageScore": 8.7,
    "highlights": "Excelente qualidade de código",
    "pointsToDevelop": "Melhorar comunicação em reuniões"
  }'
```

## 🏗️ Arquitetura

```
backend/
├── main.go              # Ponto de entrada da aplicação
├── config/             # Configurações
│   └── config.go
├── database/           # Conexão e migrações do banco
│   └── database.go
├── models/             # Modelos de dados
│   └── models.go
├── handlers/           # Handlers HTTP
│   ├── teams.go
│   ├── developers.go
│   └── performance_reports.go
└── routes/             # Definição de rotas
    └── routes.go
```

## 🔒 Recursos de Segurança

- Validação de dados de entrada
- Tratamento de erros adequado
- CORS configurado
- Sanitização de queries SQL (uso de prepared statements)

## 📈 Performance

- Conexão persistente com o banco de dados
- Índices otimizados para consultas frequentes
- Logs estruturados para debugging

## 🧪 Testing

Para executar os testes (quando implementados):
```bash
go test ./...
```

## 📦 Build para Produção

```bash
go build -o tivix-performance-tracker main.go
```
