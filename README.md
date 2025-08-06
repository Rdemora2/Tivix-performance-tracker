# Tivix Performance Tracker

Uma plataforma completa de gestão e avaliação de performance para equipes de desenvolvimento, construída com arquitetura moderna e foco na experiência do usuário.

## 🎯 Sobre o Projeto

O **Tivix Performance Tracker** é uma solução empresarial desenvolvida para facilitar o processo de avaliação de performance de desenvolvedores em ambientes corporativos. A plataforma oferece um sistema estruturado de avaliação, dashboards interativos e relatórios consolidados que permitem aos gestores acompanhar a evolução da equipe de forma eficiente e transparente.

## ✨ Características Principais

- **Sistema de Avaliação Estruturado**: Metodologia baseada em categorias ponderadas para avaliação objetiva
- **Dashboards Interativos**: Visualizações modernas com gráficos responsivos e indicadores de performance
- **Gestão Multi-empresa**: Suporte a múltiplas organizações com controle de acesso baseado em papéis
- **Relatórios Consolidados**: Geração de relatórios executivos para apresentações direcionais
- **Interface Moderna**: Design responsivo com suporte a tema claro/escuro

## 🏗️ Arquitetura

O projeto adota uma arquitetura de microserviços containerizada, garantindo escalabilidade e facilidade de deployment:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│                 │    │                 │    │                 │
│   Frontend      │◄──►│     Nginx       │◄──►│    Backend      │
│   (React/Vite)  │    │   (Proxy/LB)    │    │   (Go/Fiber)    │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                                                        ▼
                                               ┌─────────────────┐
                                               │                 │
                                               │   PostgreSQL    │
                                               │   (Database)    │
                                               │                 │
                                               └─────────────────┘
```

## 🚀 Stack Tecnológica

### Frontend

- **React 19** com Vite para build otimizado
- **Mantine UI** para componentes modernos e acessíveis
- **Zustand** para gerenciamento de estado
- **Recharts** para visualizações de dados

### Backend

- **Go** com framework Fiber para alta performance
- **PostgreSQL** com JSONB para dados estruturados
- **JWT** para autenticação e autorização
- **Docker** para containerização

### Infraestrutura

- **Nginx** como proxy reverso e load balancer
- **Docker Compose** para orquestração local
- **Traefik** ready para deploy em produção

## 📊 Funcionalidades

- **Dashboard Executivo**: Métricas consolidadas da equipe
- **Perfis Individuais**: Histórico detalhado de cada desenvolvedor
- **Sistema de Avaliação**: Formulários estruturados com categorias ponderadas
- **Relatórios Mensais**: Acompanhamento temporal da performance
- **Gestão de Usuários**: Controle de acesso com múltiplos níveis de permissão
- **Exportação de Dados**: Relatórios em PDF para apresentações

## 🔧 Deployment

O projeto está configurado para deployment em containers Docker com suporte a:

- **Desenvolvimento Local**: Docker Compose com hot reload
- **Produção**: Configurações otimizadas com proxy reverso
- **CI/CD Ready**: Estrutura preparada para pipelines automatizados

## 📝 Estrutura do Repositório

```
├── front-end/          # Aplicação React com interface moderna
├── back-end/           # API REST em Go com Fiber
├── nginx/              # Configurações de proxy e load balancing
├── docker-compose.yml  # Orquestração dos serviços
└── README.md          # Documentação principal
```

## 🎨 Interface

A interface foi desenvolvida com foco na usabilidade e experiência do usuário, apresentando:

- Design responsivo adaptável a diferentes dispositivos
- Tema claro/escuro com persistência de preferência
- Gráficos interativos para visualização de dados
- Navegação intuitiva e acessível

## 🔒 Segurança

- Autenticação JWT com refresh token
- Controle de acesso baseado em papéis (RBAC)
- Validação de dados em frontend e backend
- Configurações seguras de CORS e headers

---

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 👤 Autor

**Rdemora2**

- GitHub: [@Rdemora2](https://github.com/Rdemora2)

**Desenvolvido com excelência técnica e atenção aos detalhes para proporcionar uma experiência superior na gestão de performance de equipes de desenvolvimento.**
