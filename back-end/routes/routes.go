package routes

import (
	"github.com/gofiber/fiber/v2"
	"tivix-performance-tracker-backend/handlers"
	"tivix-performance-tracker-backend/middleware"
)

func SetupRoutes(app *fiber.App) {
	// Grupo principal da API
	api := app.Group("/api/v1")

	// Rotas públicas de autenticação
	auth := api.Group("/auth")
	auth.Post("/login", handlers.Login)

	// Rotas de inicialização do sistema
	init := api.Group("/init")
	init.Get("/check", handlers.CheckInitialization)
	init.Post("/admin", handlers.CreateAdminUser)

	// Rotas protegidas de autenticação - requerem token válido
	authProtected := api.Group("/auth", middleware.AuthMiddleware())
	authProtected.Get("/profile", handlers.GetProfile)
	authProtected.Post("/refresh", handlers.RefreshToken)
	authProtected.Post("/set-new-password", handlers.SetNewPassword)
	authProtected.Post("/change-password", handlers.ChangePassword)

	// Rotas admin e manager - para gerenciamento de usuários e empresas
	adminAndManagerAuth := authProtected.Group("/", middleware.ManagerOrAdminMiddleware())
	adminAndManagerAuth.Post("/create-user", handlers.CreateUser)
	adminAndManagerAuth.Get("/users", handlers.ListUsers)
	adminAndManagerAuth.Put("/users/:id", handlers.UpdateUser)
	adminAndManagerAuth.Delete("/users/:id", handlers.DeleteUser)
	
	// Rota para listar empresas - gerentes e admins podem acessar
	companiesListAuth := api.Group("/companies", middleware.AuthMiddleware(), middleware.ManagerOrAdminMiddleware())
	companiesListAuth.Get("/", handlers.GetAllCompanies)
	
	// Rotas admin apenas - para gerenciamento de empresas (diretamente no API, não no auth)
	companiesAdminAuth := api.Group("/companies", middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware())
	companiesAdminAuth.Post("/", handlers.CreateCompany)
	companiesAdminAuth.Get("/:id", handlers.GetCompanyByID)
	companiesAdminAuth.Put("/:id", handlers.UpdateCompany)
	companiesAdminAuth.Delete("/:id", handlers.DeleteCompany)

	// Middleware para todas as rotas protegidas - verifica se precisa trocar senha e empresa
	protectedWithPasswordCheck := api.Group("/", middleware.AuthMiddleware(), middleware.CheckPasswordChangeMiddleware(), middleware.CompanyAccessMiddleware())

	// Rotas de times - protegidas
	teams := protectedWithPasswordCheck.Group("/teams")
	teams.Get("/", handlers.GetAllTeams)
	teams.Get("/:id", handlers.GetTeamByID)
	teams.Post("/", middleware.ManagerOrAdminMiddleware(), handlers.CreateTeam)
	teams.Put("/:id", middleware.ManagerOrAdminMiddleware(), handlers.UpdateTeam)
	teams.Delete("/:id", middleware.AdminOnlyMiddleware(), handlers.DeleteTeam)

	// Rotas de desenvolvedores - protegidas
	developers := protectedWithPasswordCheck.Group("/developers")
	developers.Get("/", handlers.GetAllDevelopers)
	developers.Get("/archived", handlers.GetArchivedDevelopers)
	developers.Get("/:id", handlers.GetDeveloperByID)
	developers.Post("/", middleware.ManagerOrAdminMiddleware(), handlers.CreateDeveloper)
	developers.Put("/:id", middleware.ManagerOrAdminMiddleware(), handlers.UpdateDeveloper)
	developers.Put("/:id/archive", middleware.ManagerOrAdminMiddleware(), handlers.ArchiveDeveloper)
	developers.Delete("/:id", middleware.ManagerOrAdminMiddleware(), handlers.DeleteDeveloper)

	// Rotas de desenvolvedores por time - protegidas
	teams.Get("/:teamId/developers", handlers.GetDevelopersByTeam)

	// Rotas de relatórios de performance - protegidas
	reports := protectedWithPasswordCheck.Group("/performance-reports")
	reports.Get("/", handlers.GetAllPerformanceReports)
	reports.Get("/months", handlers.GetAvailableMonths)
	reports.Get("/stats", handlers.GetPerformanceStats)
	reports.Get("/:id", handlers.GetPerformanceReportByID)
	reports.Post("/", middleware.ManagerOrAdminMiddleware(), handlers.CreatePerformanceReport)

	// Rotas de relatórios por desenvolvedor - protegidas
	developers.Get("/:developerId/reports", handlers.GetPerformanceReportsByDeveloper)

	// Rotas de relatórios por mês - protegidas
	reports.Get("/month/:month", handlers.GetPerformanceReportsByMonth)
}
