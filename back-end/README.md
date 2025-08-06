# Backend - Tivix Performance Tracker

API REST de alta performance construída em Go com arquitetura escalável para gestão de avaliações de performance de equipes de desenvolvimento.

## 🎯 Visão Técnica

Esta API representa uma implementação moderna de microserviços em Go, utilizando o framework Fiber para máxima performance e PostgreSQL como banco de dados principal. A arquitetura foi projetada para suportar cargas elevadas, múltiplas empresas e diferentes níveis de acesso com segurança enterprise-grade.

## 🏗️ Arquitetura de Software

A aplicação segue o padrão **Clean**

## 🚀 Execução Local

### Desenvolvimento\*\* com separação clara de responsabilidades:

```
cmd/                    # Entry points da aplicação
├── migration-status/  # Utilitário para verificar status das migrações
config/                # Configurações e variáveis de ambiente
database/              # Conexão e execução de migrações
handlers/              # Controllers/Handlers HTTP
├── auth.go           # Autenticação e autorização
├── companies.go      # Gestão de empresas
├── developers.go     # CRUD de desenvolvedores
├── performance_reports.go # Core business logic
└── teams.go          # Gestão de equipes
middleware/            # Middlewares de autenticação e autorização
migrations/            # Sistema centralizado de migrações
├── README.md         # Documentação das migrações
├── manager.go        # Gerenciador de migrações
├── sql_migrations.go # Definições SQL das migrações
└── *.sql            # Arquivos individuais de migração
models/               # Entidades de domínio e DTOs
routes/               # Definição de rotas e agrupamentos
utils/                # Utilitários e helpers
```

### Padrões de Design Implementados

- **Repository Pattern**: Abstração da camada de dados
- **Middleware Pattern**: Cross-cutting concerns (auth, logging, CORS)
- **DTO Pattern**: Data Transfer Objects para API contracts
- **Factory Pattern**: Criação de objetos complexos
- **Dependency Injection**: Inversão de controle para testabilidade

## 🛠️ Stack Tecnológica Detalhada

### Core Framework

- **Go 1.24**: Linguagem de alta performance com garbage collector otimizado
- **Fiber v2**: Framework web inspirado no Express.js, extremamente rápido
- **SQLX**: Driver PostgreSQL com suporte a SQL raw e mapeamento automático

### Banco de Dados

- **PostgreSQL 15+**: RDBMS robusto com suporte JSONB nativo
- **UUID Extensions**: Identificadores únicos universais
- **JSONB**: Armazenamento flexível para dados semi-estruturados

### Segurança e Autenticação

- **JWT (golang-jwt/jwt/v5)**: Tokens stateless com refresh automático
- **bcrypt**: Hash de senhas com salt automático
- **CORS**: Configuração granular de Cross-Origin Resource Sharing

### Validação e Configuração

- **go-playground/validator**: Validação estrutural de dados
- **godotenv**: Gerenciamento de variáveis de ambiente
- **UUID (google/uuid)**: Geração e manipulação de identificadores únicos

## 🗄️ Modelagem de Dados

### Esquema Relacional

```sql
-- Empresas (Multi-tenancy)
companies (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  is_active BOOLEAN DEFAULT true,
  created_at, updated_at TIMESTAMP
);

-- Usuários com RBAC
users (
  id UUID PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  role ENUM('admin', 'manager', 'user'),
  company_id UUID REFERENCES companies(id),
  needs_password_change BOOLEAN DEFAULT false,
  is_active BOOLEAN DEFAULT true,
  created_at, updated_at TIMESTAMP
);

-- Equipes organizacionais
teams (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  color VARCHAR(50) DEFAULT 'blue',
  company_id UUID REFERENCES companies(id),
  created_at, updated_at TIMESTAMP
);

-- Desenvolvedores
developers (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  role VARCHAR(255) NOT NULL,
  latest_performance_score DECIMAL(4,2) DEFAULT 0.00,
  team_id UUID REFERENCES teams(id),
  company_id UUID REFERENCES companies(id),
  archived_at TIMESTAMP NULL,
  created_at, updated_at TIMESTAMP
);

-- Relatórios de Performance (Core Business)
performance_reports (
  id UUID PRIMARY KEY,
  developer_id UUID REFERENCES developers(id),
  month VARCHAR(7) NOT NULL, -- YYYY-MM format
  question_scores JSONB NOT NULL,
  category_scores JSONB NOT NULL,
  weighted_average_score DECIMAL(4,2) NOT NULL,
  highlights TEXT,
  points_to_develop TEXT,
  created_at, updated_at TIMESTAMP
);
```

### Índices de Performance

```sql
-- Índices estratégicos para queries frequentes
CREATE INDEX idx_performance_reports_month ON performance_reports(month);
CREATE INDEX idx_performance_reports_developer ON performance_reports(developer_id);
CREATE INDEX idx_performance_reports_score ON performance_reports(weighted_average_score);
CREATE INDEX idx_developers_company ON developers(company_id);
CREATE INDEX idx_users_company_role ON users(company_id, role);
```

## 🔒 Sistema de Segurança

### Autenticação JWT

```go
// Estrutura do JWT Claims
type JWTClaims struct {
    UserID    uuid.UUID  `json:"user_id"`
    Email     string     `json:"email"`
    Role      string     `json:"role"`
    CompanyID *uuid.UUID `json:"company_id,omitempty"`
    jwt.RegisteredClaims
}

// Middleware de autenticação com refresh automático
func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := extractToken(c)
        claims, err := validateJWT(token)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "error": "Invalid or expired token",
            })
        }
        c.Locals("user", claims)
        return c.Next()
    }
}
```

### Controle de Acesso Baseado em Papéis (RBAC)

```go
// Níveis de permissão hierárquicos
const (
    RoleAdmin   = "admin"    // Acesso total ao sistema
    RoleManager = "manager"  // Gestão dentro da empresa
    RoleUser    = "user"     // Acesso limitado de visualização da empresa
)

// Middleware de autorização por papel
func AdminOnlyMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        user := c.Locals("user").(*JWTClaims)
        if user.Role != RoleAdmin {
            return c.Status(403).JSON(fiber.Map{
                "error": "Insufficient permissions",
            })
        }
        return c.Next()
    }
}
```

### Multi-tenancy (Isolamento por Empresa)

```go
// Middleware de isolamento de dados por empresa
func CompanyAccessMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        user := c.Locals("user").(*JWTClaims)

        // Admins podem acessar qualquer empresa
        if user.Role == RoleAdmin {
            return c.Next()
        }

        // Usuários devem ter empresa associada
        if user.CompanyID == nil {
            return c.Status(403).JSON(fiber.Map{
                "error": "User must be associated with a company",
            })
        }

        return c.Next()
    }
}
```

## 🌐 API Design e Endpoints

### Estrutura RESTful

```
/api/v1/
├── auth/                          # Autenticação
│   ├── POST /login               # Login com email/password
│   ├── GET /profile              # Perfil do usuário logado
│   ├── POST /refresh             # Refresh do JWT token
│   └── POST /set-new-password    # Alteração de senha obrigatória
├── init/                         # Inicialização do sistema
│   ├── GET /check               # Verificar se sistema foi inicializado
│   └── POST /admin              # Criar primeiro usuário admin
├── companies/                    # Gestão de empresas (Admin only)
│   ├── GET /                    # Listar empresas
│   ├── POST /                   # Criar empresa
│   ├── GET /:id                 # Detalhes da empresa
│   ├── PUT /:id                 # Atualizar empresa
│   └── DELETE /:id              # Remover empresa
├── teams/                       # Gestão de equipes
│   ├── GET /                    # Listar equipes da empresa
│   ├── POST /                   # Criar equipe
│   ├── PUT /:id                 # Atualizar equipe
│   └── DELETE /:id              # Remover equipe
├── developers/                  # CRUD de desenvolvedores
│   ├── GET /                    # Listar desenvolvedores
│   ├── POST /                   # Adicionar desenvolvedor
│   ├── GET /:id                 # Detalhes do desenvolvedor
│   ├── PUT /:id                 # Atualizar desenvolvedor
│   ├── DELETE /:id              # Arquivar desenvolvedor
│   └── POST /:id/restore        # Restaurar desenvolvedor
└── performance-reports/         # Core business - Relatórios
    ├── GET /                    # Listar todos os relatórios
    ├── POST /                   # Criar novo relatório
    ├── GET /:id                 # Detalhes de relatório específico
    ├── GET /developer/:id       # Relatórios por desenvolvedor
    ├── GET /month/:month        # Relatórios por mês
    ├── GET /months              # Meses com relatórios disponíveis
    └── GET /stats               # Estatísticas consolidadas
```

### Padronização de Responses

```go
// Resposta de sucesso padrão
type SuccessResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data"`
    Message string      `json:"message,omitempty"`
}

// Resposta de erro padrão
type ErrorResponse struct {
    Error   bool   `json:"error"`
    Message string `json:"message"`
    Code    string `json:"code,omitempty"`
}
```

## 📊 Business Logic - Sistema de Performance

### Algoritmo de Cálculo de Performance

```go
// Estrutura das categorias de avaliação
type EvaluationCategory struct {
    Label     string    `json:"label"`
    Weight    float64   `json:"weight"`
    Questions []Question `json:"questions"`
}

type Question struct {
    Key    string  `json:"key"`
    Label  string  `json:"label"`
    Weight float64 `json:"weight"`
}

// Categorias definidas no sistema
var EvaluationCategories = map[string]EvaluationCategory{
    "commitment": {
        Label:  "Comprometimento e Disciplina",
        Weight: 0.3,
        Questions: []Question{
            {Key: "punctualityDeliveries", Label: "Pontualidade nas Entregas", Weight: 3},
            {Key: "punctualityRituals", Label: "Pontualidade em Rituais", Weight: 2},
            {Key: "hybridModelAdherence", Label: "Adesão ao Modelo Híbrido", Weight: 1},
        },
    },
    "technicalQuality": {
        Label:  "Qualidade e Execução Técnica",
        Weight: 0.4,
        Questions: []Question{
            {Key: "deliveryQuality", Label: "Qualidade das Entregas", Weight: 4},
            {Key: "taskAutonomy", Label: "Autonomia na Resolução de Tarefas", Weight: 3},
        },
    },
    "collaboration": {
        Label:  "Colaboração e Proatividade",
        Weight: 0.3,
        Questions: []Question{
            {Key: "proactivityImprovements", Label: "Proatividade e Sugestão de Melhorias", Weight: 3},
            {Key: "communicationQuality", Label: "Qualidade da Comunicação", Weight: 2},
        },
    },
}
```

### Validação e Processamento de Relatórios

```go
// Handler para criação de relatório de performance
func CreatePerformanceReport(c *fiber.Ctx) error {
    var req models.CreatePerformanceReportRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(ErrorResponse{
            Error:   true,
            Message: "Invalid request format",
        })
    }

    // Validações de negócio
    if err := validatePerformanceReport(&req); err != nil {
        return c.Status(400).JSON(ErrorResponse{
            Error:   true,
            Message: err.Error(),
        })
    }

    // Verificar duplicidade de relatório no mês
    if exists := checkExistingReport(req.DeveloperID, req.Month); exists {
        return c.Status(409).JSON(ErrorResponse{
            Error:   true,
            Message: "Report already exists for this month",
        })
    }

    // Salvar no banco com transação
    report, err := savePerformanceReport(&req)
    if err != nil {
        return c.Status(500).JSON(ErrorResponse{
            Error:   true,
            Message: "Failed to save report",
        })
    }

    // Atualizar score mais recente do desenvolvedor
    updateDeveloperLatestScore(req.DeveloperID, req.WeightedAverageScore)

    return c.Status(201).JSON(SuccessResponse{
        Success: true,
        Data:    report,
        Message: "Performance report created successfully",
    })
}
```

## ⚡ Otimizações de Performance

### Database Query Optimization

```go
// Query otimizada para relatórios consolidados
func GetConsolidatedReports(companyID uuid.UUID, month string) ([]PerformanceReport, error) {
    query := `
        SELECT
            pr.id, pr.developer_id, pr.month, pr.question_scores,
            pr.category_scores, pr.weighted_average_score,
            pr.highlights, pr.points_to_develop, pr.created_at, pr.updated_at,
            d.name as developer_name, d.role as developer_role
        FROM performance_reports pr
        INNER JOIN developers d ON pr.developer_id = d.id
        WHERE d.company_id = $1 AND pr.month = $2
        ORDER BY pr.weighted_average_score DESC, pr.created_at DESC
    `

    var reports []PerformanceReport
    err := database.DB.Select(&reports, query, companyID, month)
    return reports, err
}
```

### Connection Pooling e Configurações

```go
// Configuração otimizada do pool de conexões
func configureDB(db *sqlx.DB) {
    db.SetMaxOpenConns(25)                 // Máximo de conexões abertas
    db.SetMaxIdleConns(5)                  // Conexões idle no pool
    db.SetConnMaxLifetime(30 * time.Minute) // Tempo de vida da conexão
    db.SetConnMaxIdleTime(5 * time.Minute)  // Tempo máximo idle
}
```

### Middleware de Logging Estruturado

```go
// Logger configurado para produção
app.Use(logger.New(logger.Config{
    Format: "[${time}] ${status} - ${method} ${path} - ${ip} - ${latency}\n",
    TimeFormat: "2006-01-02 15:04:05",
    TimeZone: "UTC",
}))
```

## 🔧 Configuração e Environment

### Configuration Management

```go
type Config struct {
    // Database
    DBHost     string `env:"DB_HOST" envDefault:"localhost"`
    DBPort     string `env:"DB_PORT" envDefault:"5432"`
    DBUser     string `env:"DB_USER" envDefault:"postgres"`
    DBPassword string `env:"DB_PASSWORD" envDefault:"postgres"`
    DBName     string `env:"DB_NAME" envDefault:"tivix_performance_tracker"`
    DBSSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`

    // Server
    Port        string `env:"PORT" envDefault:"8080"`
    Host        string `env:"HOST" envDefault:"localhost"`
    Environment string `env:"ENVIRONMENT" envDefault:"development"`

    // Security
    JWTSecret  string `env:"JWT_SECRET" envDefault:"change-in-production"`
    CORSOrigin string `env:"CORS_ORIGIN" envDefault:"http://localhost:5173"`
}
```

### Docker Configuration

```dockerfile
# Multi-stage build para otimização
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

## 🚀 Deploy e DevOps

### Health Checks

```go
// Endpoint de health check com verificações de dependências
app.Get("/health", func(c *fiber.Ctx) error {
    // Verificar conexão com banco
    if err := database.DB.Ping(); err != nil {
        return c.Status(503).JSON(fiber.Map{
            "status": "unhealthy",
            "checks": map[string]string{
                "database": "failing",
            },
        })
    }

    return c.JSON(fiber.Map{
        "status": "healthy",
        "timestamp": time.Now().UTC(),
        "version": "1.0.0",
        "checks": map[string]string{
            "database": "passing",
        },
    })
})
```

### Graceful Shutdown

```go
// Implementação de shutdown graceful
func gracefulShutdown(app *fiber.App) {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)

    go func() {
        <-c
        log.Println("Gracefully shutting down...")
        app.Shutdown()
    }()

    if err := app.Listen(":8080"); err != nil {
        log.Panic(err)
    }
}
```

## 📈 Monitoramento e Observabilidade

### Métricas de Performance

- **Latência**: P50, P95, P99 de response time
- **Throughput**: Requests por segundo
- **Error Rate**: Percentage de erros 4xx/5xx
- **Database**: Connection pool utilization, query duration

### Structured Logging

```go
// Log estruturado para observabilidade
log.WithFields(logrus.Fields{
    "user_id":    userID,
    "company_id": companyID,
    "operation":  "create_performance_report",
    "duration":   duration,
}).Info("Performance report created successfully")
```

## �️ Sistema de Migrações

### Visão Geral

O sistema utiliza um gerenciador de migrações centralizado que garante:

- **Versionamento sequencial** das mudanças no banco
- **Controle de estado** com tabela `schema_migrations`
- **Execução transacional** para rollback automático em caso de erro
- **Documentação completa** de cada migração

### Estrutura das Migrações

```text
migrations/
├── README.md                     # Documentação completa
├── manager.go                    # Gerenciador de migrações
├── sql_migrations.go             # Definições SQL embeddadas
├── 001_initial_setup.sql         # Configuração PostgreSQL
├── 002_create_tables.sql         # Tabelas principais
├── 003_create_indexes.sql        # Índices de performance
├── 004_create_triggers.sql       # Triggers para timestamps
├── 005_multitenant_implementation.sql # Sistema multitenant
└── 006_data_migration_multitenant.sql # Migração de dados
```

### Execução Automática

As migrações são executadas **automaticamente** quando a aplicação inicia:

```bash
# Executar aplicação (migrações automáticas)
go run main.go

# Verificar status das migrações
go run cmd/migration-status/main.go
```

### Exemplo de Output

```text
📊 Verificando status das migrações...

ID                              Descrição                           Status        Data
---                             ----------                          ------        ----
001_initial_setup               Configuração inicial PostgreSQL    ✅ Aplicada   2025-08-05 10:30:15
002_create_tables               Criação das tabelas principais      ✅ Aplicada   2025-08-05 10:30:16
003_create_indexes              Criação de índices para performance ✅ Aplicada   2025-08-05 10:30:17
004_create_triggers             Configuração de triggers            ✅ Aplicada   2025-08-05 10:30:18
005_multitenant_implementation  Implementação do sistema multitenant ⏳ Pendente  -
006_data_migration_multitenant  Migração de dados para multitenancy  ⏳ Pendente  -

📈 Resumo das Migrações:
   • Total: 6
   • Aplicadas: 4
   • Pendentes: 2
```

### Criando Nova Migração

1. **Adicionar arquivo SQL** na pasta `migrations/`
2. **Atualizar `sql_migrations.go`** com a nova constante
3. **Adicionar à lista** em `manager.go` no método `GetAllMigrations()`
4. **Documentar** no `README.md` das migrações

Exemplo:

```go
// Em sql_migrations.go
const migration007SQL = `
-- Sua nova migração aqui
ALTER TABLE users ADD COLUMN last_login TIMESTAMP;
`

// Em manager.go
{
    ID:          "007_add_last_login",
    Description: "Adicionar campo last_login na tabela users",
    SQL:         migration007SQL,
},
```

## �🚀 Execução Local

### Desenvolvimento

```bash
# Instalar dependências
go mod download

# Executar aplicação (migrações automáticas)
go run main.go

# Executar com hot reload (air)
air
```

### Configuração de Ambiente

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=db_name
DB_SSLMODE=disable

# Server Configuration
PORT=8080
HOST=localhost
ENVIRONMENT=development

# Security
JWT_SECRET=your-secret-key-change-in-production
CORS_ORIGIN=http://localhost:5173
```

### Build para Produção

```bash
# Build otimizado
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Ou usando Docker
docker build -t tivix-backend .
```

---

**Esta API representa uma implementação enterprise-grade em Go, priorizando performance, segurança e escalabilidade para ambientes de produção críticos.**
