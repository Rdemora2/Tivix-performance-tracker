# Frontend - Tivix Performance Tracker

Interface moderna e responsiva para gestão de performance de equipes de desenvolvimento, construída com as mais recentes tecnologias React e focada na experiência do usuário.

## 🎯 Visão Técnica

Este frontend representa uma aplicação **Single Page Application (SPA)** de alta performance, desenvolvida com React 19 e arquitetura de componentes modulares. A aplicação implementa padrões modernos de desenvolvimento, incluindo hooks customizados, gerenciamento de estado reativo e design system consistente.

## 🏗️ Arquitetura da Aplicação

### Estrutura Modular

```
src/
├── components/          # Componentes reutilizáveis
│   ├── ui/             # Design System (Radix UI + Mantine)
│   ├── ProtectedRoute/ # Controle de acesso e autenticação
│   └── PermissionGuard/ # Autorização baseada em papéis
├── hooks/              # Hooks customizados para lógica compartilhada
├── layouts/            # Templates de layout da aplicação
├── pages/              # Componentes de página (route-based)
├── services/           # Camada de comunicação com APIs
├── store/              # Gerenciamento de estado global (Zustand)
├── types/              # Definições de tipos e modelos de dados
└── utils/              # Funções utilitárias e helpers
```

### Padrões de Design Implementados

- **Container/Presentational**: Separação clara entre lógica de negócio e apresentação
- **Custom Hooks**: Encapsulamento de lógica stateful reutilizável
- **Compound Components**: Componentes complexos com API declarativa
- **Render Props**: Compartilhamento de lógica entre componentes

## 🛠️ Stack Tecnológica Detalhada

### Core Framework

- **React 19**: Utilização das mais recentes funcionalidades, incluindo Concurrent Features
- **React Router DOM v7**: Roteamento declarativo com lazy loading de rotas
- **Vite 6**: Build tool otimizado com HMR (Hot Module Replacement) e bundling eficiente

### UI/UX Libraries

- **Mantine UI 8**: Sistema de componentes moderno com tema customizável
- **Radix UI**: Primitivos acessíveis para componentes complexos (dialogs, dropdowns, etc.)
- **Tabler Icons**: Biblioteca consistente de ícones vetoriais
- **Framer Motion**: Animações performáticas e transições fluidas

### Visualização de Dados

- **Recharts**: Gráficos responsivos e interativos com D3.js
- **HTML2Canvas + jsPDF**: Geração de relatórios em PDF do lado cliente

### Gerenciamento de Estado

- **Zustand**: Store reativo e minimalista para estado global
- **React Hook Form**: Gerenciamento otimizado de formulários com validação
- **Zod**: Schema validation com TypeScript inference

### Estilização e Design System

- **Tailwind CSS v4**: Utility-first CSS framework para prototipagem rápida
- **CSS Variables**: Sistema de temas dinâmico (claro/escuro)
- **Responsive Design**: Mobile-first approach com breakpoints otimizados

## 🎨 Sistema de Design

### Paleta de Cores Estratégica

```css
/* Performance Indicators */
--success: #51cf66      /* Performance Excelente (≥8.0) */
--primary: #228be6      /* Performance Boa (6.0-7.9) */
--warning: #ffd43b      /* Performance Regular (4.0-5.9) */
--danger: #ff6b6b       /* Necessita Melhoria (<4.0) */
```

### Typography Scale

- **Fonte Principal**: Inter - otimizada para legibilidade em interfaces
- **Hierarchy**: Sistema de escalas baseado em modular scale (1.250 - Major Third)

### Component Design Tokens

- **Spacing**: Sistema de 8px grid para consistência visual
- **Border Radius**: Valores padronizados (4px, 8px, 12px)
- **Shadows**: Elevações sutis para hierarquia de informação

## 📊 Funcionalidades Técnicas Implementadas

### Sistema de Avaliação Avançado

```javascript
// Algoritmo de cálculo de performance ponderada
const calculatePerformance = (questionScores, categories) => {
  return Object.entries(categories).reduce((total, [key, category]) => {
    const categoryScore =
      category.questions.reduce((sum, question) => {
        return sum + questionScores[question.key] * question.weight;
      }, 0) / category.questions.reduce((sum, q) => sum + q.weight, 0);

    return total + categoryScore * category.weight;
  }, 0);
};
```

### Arquitetura de Estado Reativo

```javascript
// Store Zustand com middleware de persistência
const useAppStore = create(
  persist(
    (set, get) => ({
      // Estado da aplicação
      developers: [],
      performanceReports: [],
      teams: [],

      // Actions com otimistic updates
      addPerformanceReport: async (reportData) => {
        // Implementação com error handling e rollback
      },
    }),
    { name: "tivix-performance-storage" }
  )
);
```

### Sistema de Roteamento Inteligente

- **Code Splitting**: Lazy loading automático de rotas
- **Protected Routes**: Autenticação e autorização em nível de rota
- **Breadcrumb Generation**: Navegação contextual automática

## 🔒 Implementação de Segurança

### Autenticação e Autorização

- **JWT Token Management**: Refresh automático com interceptors
- **Role-Based Access Control (RBAC)**: Controle granular de permissões
- **XSS Protection**: Sanitização de inputs e CSP headers

### Validação de Dados

```javascript
// Schema Zod para validação robusta
const performanceReportSchema = z.object({
  developerId: z.string().uuid(),
  month: z.string().regex(/^\d{4}-\d{2}$/),
  questionScores: z.record(z.number().min(0).max(10)),
  weightedAverageScore: z.number().min(0).max(10),
  highlights: z.string().optional(),
  pointsToDevelop: z.string().optional(),
});
```

## 📱 Responsividade e Performance

### Mobile-First Approach

- **Breakpoints**: Customizados para diferentes dispositivos
- **Touch Optimization**: Gestos e interações otimizadas para mobile
- **Progressive Enhancement**: Funcionalidades básicas sempre disponíveis

### Otimizações de Performance

- **Bundle Splitting**: Chunks otimizados por rota
- **Image Optimization**: Lazy loading e formatos modernos
- **Memoization**: React.memo e useMemo estratégicos
- **Virtual Scrolling**: Para listas extensas de dados

## 🧪 Qualidade de Código

### Padrões Implementados

- **ESLint**: Configuração strict com regras customizadas
- **Prettier**: Formatação automática consistente
- **Husky**: Git hooks para validação pré-commit

### Estrutura de Componentes

```jsx
// Padrão de component composition
const PerformanceCard = ({ developer, onViewProfile, onCreateReport }) => {
  const { latestScore, trend } = usePerformanceAnalytics(developer.id);

  return (
    <Card>
      <PerformanceIndicator score={latestScore} trend={trend} />
      <DeveloperInfo developer={developer} />
      <ActionGroup>
        <Button onClick={onViewProfile}>Ver Perfil</Button>
        <Button onClick={onCreateReport}>Nova Avaliação</Button>
      </ActionGroup>
    </Card>
  );
};
```

## 🚀 Build e Deploy

### Desenvolvimento Local

```bash
# Instalação com npm para performance otimizada
npm install

# Desenvolvimento com hot reload
npm run dev

# Build para produção com otimizações
npm run build
```

### Configuração de Build

- **Vite Plugins**: React, Tailwind CSS, Bundle Analyzer
- **Environment Variables**: Configuração por ambiente
- **Tree Shaking**: Eliminação de código não utilizado
- **Asset Optimization**: Compressão e otimização automática

## 📈 Métricas e Analytics

### Performance Monitoring

- **Web Vitals**: Monitoramento de métricas de performance
- **Bundle Size Tracking**: Análise de tamanho de bundles
- **Render Performance**: Profiling de componentes React

### User Experience

- **Loading States**: Feedback visual em todas as interações
- **Error Boundaries**: Tratamento graceful de erros
- **Accessibility**: WCAG 2.1 compliance

## 🔧 Configurações Avançadas

### Vite Configuration

```javascript
export default defineConfig({
  plugins: [react(), tailwindcss()],
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ["react", "react-dom"],
          ui: ["@mantine/core", "@radix-ui/react-dialog"],
          charts: ["recharts"],
        },
      },
    },
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
});
```

### PWA Ready

- **Service Worker**: Cache strategy configurada
- **Offline Support**: Funcionalidades básicas offline
- **Install Prompt**: App installation nativa

---

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
- npm

### Instalação e Execução

```bash
# Navegar para o diretório do projeto
cd tivix-performance-tracker

# Instalar dependências
npm install

# Executar em modo de desenvolvimento
npm run dev

# Build para produção
npm run build
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
