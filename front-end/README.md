# Performance Tracker

Uma aplicação web moderna para avaliação e acompanhamento da performance de equipes de desenvolvimento.

## 🚀 Visão Geral

O **Performance Tracker** é uma ferramenta completa que permite aos gestores avaliar a performance mensal de suas equipes de desenvolvedores de forma estruturada, intuitiva e visualmente atraente.

## ✨ Funcionalidades Principais

### 📊 Dashboard da Equipe

- Visualização em cards de todos os membros da equipe
- Exibição da performance atual de cada desenvolvedor
- Badges visuais indicando o nível de performance (Excelente, Bom, Regular, Precisa Melhorar)
- Acesso rápido para criar novas avaliações ou visualizar perfis

### 👤 Perfil Individual do Desenvolvedor

- Página dedicada para cada membro da equipe
- Gráfico de linha mostrando a evolução da performance ao longo do tempo
- Histórico completo de relatórios de avaliação
- Gráfico de radar para visualizar pontos fortes e fracos em avaliações específicas

### 📝 Sistema de Avaliação Estruturado

- Formulário dividido em etapas lógicas usando Stepper
- Avaliação por categorias com sliders interativos (0-10):
  - **Comprometimento** (peso: 25%)
  - **Qualidade Técnica** (peso: 35%)
  - **Colaboração** (peso: 20%)
  - **Resolução de Problemas** (peso: 20%)
- Cálculo automático da nota final usando média ponderada
- Campos para feedback qualitativo (Destaques e Pontos a Desenvolver)

### 📈 Relatório Consolidado

- Visão geral da performance de toda a equipe por período
- Estatísticas consolidadas (média, maior e menor nota)
- Tabela ordenável com todas as avaliações do período
- Seção de observações qualitativas
- Funcionalidade de exportação para PDF (estrutura implementada)

### 🌓 Modo Claro/Escuro

- Alternância completa entre temas claro e escuro
- Persistência da preferência do usuário
- Adaptação automática de todos os componentes e gráficos

## 🛠️ Tecnologias Utilizadas

### Frontend

- **React 19** - Biblioteca principal
- **Mantine UI 8** - Biblioteca de componentes moderna
- **Recharts** - Gráficos interativos e responsivos
- **React Router DOM** - Navegação entre páginas
- **Zustand** - Gerenciamento de estado simples e eficiente
- **Tabler Icons** - Ícones consistentes e modernos

### Ferramentas de Desenvolvimento

- **Vite** - Build tool rápido e moderno
- **JavaScript (JSX)** - Linguagem principal
- **CSS Modules** - Estilização modular

## 🏗️ Arquitetura

A aplicação segue uma arquitetura baseada em componentes com separação clara de responsabilidades:

```
src/
├── components/     # Componentes reutilizáveis
├── pages/         # Páginas da aplicação
├── layouts/       # Layouts de página
├── store/         # Gerenciamento de estado (Zustand)
├── types/         # Definições de tipos e interfaces
├── utils/         # Funções utilitárias
└── services/      # Lógica de negócio e APIs
```

## 📋 Como Executar

### Pré-requisitos

- Node.js 20+
- pnpm (recomendado) ou npm

### Instalação e Execução

```bash
# Navegar para o diretório do projeto
cd tivix-performance-tracker

# Instalar dependências
pnpm install

# Executar em modo de desenvolvimento
pnpm run dev

# Build para produção
pnpm run build
```

A aplicação estará disponível em `http://localhost:5173`

## 🎯 Fluxo de Uso

1. **Dashboard**: Visualize todos os membros da equipe e suas performances atuais
2. **Nova Avaliação**: Clique em "Nova Avaliação" para criar um relatório mensal
3. **Preenchimento**: Complete as 3 etapas do formulário (Período, Avaliações, Comentários)
4. **Perfil Individual**: Acesse o perfil para ver histórico e evolução
5. **Relatório Consolidado**: Gere relatórios para apresentação à diretoria

## 🎨 Design e UX

### Princípios de Design

- **Minimalismo**: Interface limpa e focada no essencial
- **Clareza**: Informações organizadas de forma intuitiva
- **Consistência**: Uso padronizado de cores, tipografia e espaçamentos
- **Acessibilidade**: Suporte a modo escuro e navegação por teclado

### Paleta de Cores

- **Primária**: Azul (#228be6) - Confiança e profissionalismo
- **Sucesso**: Verde - Performance excelente
- **Atenção**: Amarelo/Laranja - Performance regular
- **Alerta**: Vermelho - Performance que precisa melhorar

## 📊 Sistema de Pontuação

### Categorias de Avaliação

- **Comprometimento** (25%): Dedicação, pontualidade, proatividade
- **Qualidade Técnica** (35%): Código limpo, boas práticas, conhecimento técnico
- **Colaboração** (20%): Trabalho em equipe, comunicação, mentoria
- **Resolução de Problemas** (20%): Análise, criatividade, eficiência

### Classificação de Performance

- **8.0 - 10.0**: Excelente (Verde)
- **6.0 - 7.9**: Bom (Azul)
- **4.0 - 5.9**: Regular (Amarelo)
- **0.0 - 3.9**: Precisa Melhorar (Vermelho)

## 🔧 Funcionalidades Técnicas

### Persistência de Dados

- Armazenamento local usando localStorage
- Sincronização automática entre abas
- Backup automático das avaliações

### Responsividade

- Design adaptativo para desktop, tablet e mobile
- Componentes otimizados para touch
- Navegação otimizada para diferentes tamanhos de tela

### Performance

- Lazy loading de componentes
- Otimização de re-renders
- Bundle splitting automático

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 👤 Autor

**Rdemora2**

- GitHub: [@Rdemora2](https://github.com/Rdemora2)
