package middleware

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"tivix-performance-tracker-backend/models"
)

type JWTClaims struct {
	UserID              uuid.UUID `json:"userId"`
	Email               string    `json:"email"`
	Role                string    `json:"role"`
	IsActive            bool      `json:"isActive"`
	NeedsPasswordChange bool      `json:"needsPasswordChange"`
	jwt.RegisteredClaims
}

func GenerateJWT(user models.User) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-this-in-production"
	}

	claims := JWTClaims{
		UserID:              user.ID,
		Email:               user.Email,
		Role:                user.Role,
		IsActive:            user.IsActive,
		NeedsPasswordChange: user.NeedsPasswordChange,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expira em 24 horas
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "tivix-performance-tracker",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ValidateJWT(tokenString string) (*JWTClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-this-in-production"
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Token de autorização não fornecido",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Formato de token inválido. Use 'Bearer <token>'",
			})
		}

		claims, err := ValidateJWT(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Token inválido ou expirado",
			})
		}

		if !claims.IsActive {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": "Usuário inativo",
			})
		}

		c.Locals("user", claims)

		return c.Next()
	}
}

func ManagerOrAdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*JWTClaims)
		if user.Role != "admin" && user.Role != "manager" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Acesso negado. Apenas administradores e gerentes têm permissão",
			})
		}
		return c.Next()
	}
}

func AdminOnlyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*JWTClaims)
		if user.Role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Acesso negado. Apenas administradores têm permissão",
			})
		}
		return c.Next()
	}
}

func CheckPasswordChangeMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		skipRoutes := []string{
			"/api/v1/auth/set-new-password",
			"/api/v1/auth/logout",
			"/api/v1/auth/profile",
			"/api/v1/auth/change-password",
			"/api/v1/auth/refresh",
		}

		path := c.Path()
		for _, route := range skipRoutes {
			if path == route {
				return c.Next()
			}
		}

		user := c.Locals("user").(*JWTClaims)
		if user.NeedsPasswordChange {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":                 "error",
				"message":                "Você deve definir uma nova senha antes de continuar",
				"requiresPasswordChange": true,
			})
		}

		return c.Next()
	}
}
