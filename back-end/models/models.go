package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONB)
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j)
	case string:
		return json.Unmarshal([]byte(v), j)
	default:
		return errors.New("cannot scan into JSONB")
	}
}

type Team struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Color       string    `json:"color" db:"color"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type Developer struct {
	ID                     uuid.UUID  `json:"id" db:"id"`
	Name                   string     `json:"name" db:"name"`
	Role                   string     `json:"role" db:"role"`
	LatestPerformanceScore float64    `json:"latestPerformanceScore" db:"latest_performance_score"`
	TeamID                 *uuid.UUID `json:"teamId" db:"team_id"`
	ArchivedAt             *time.Time `json:"archivedAt" db:"archived_at"`
	CreatedAt              time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt              time.Time  `json:"updatedAt" db:"updated_at"`
}

type PerformanceReport struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	DeveloperID          uuid.UUID `json:"developerId" db:"developer_id"`
	Month                string    `json:"month" db:"month"`
	QuestionScores       JSONB     `json:"questionScores" db:"question_scores"`
	CategoryScores       JSONB     `json:"categoryScores" db:"category_scores"`
	WeightedAverageScore float64   `json:"weightedAverageScore" db:"weighted_average_score"`
	Highlights           string    `json:"highlights" db:"highlights"`
	PointsToDevelop      string    `json:"pointsToDevelop" db:"points_to_develop"`
	CreatedAt            time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt            time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateTeamRequest struct {
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type UpdateTeamRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
}

type CreateDeveloperRequest struct {
	Name   string     `json:"name" validate:"required,min=2"`
	Role   string     `json:"role" validate:"required,min=2"`
	TeamID *uuid.UUID `json:"teamId,omitempty"`
}

type UpdateDeveloperRequest struct {
	Name                   *string    `json:"name,omitempty"`
	Role                   *string    `json:"role,omitempty"`
	LatestPerformanceScore *float64   `json:"latestPerformanceScore,omitempty"`
	TeamID                 *uuid.UUID `json:"teamId,omitempty"`
}

type CreatePerformanceReportRequest struct {
	DeveloperID          uuid.UUID `json:"developerId" validate:"required"`
	Month                string    `json:"month" validate:"required"`
	QuestionScores       JSONB     `json:"questionScores" validate:"required"`
	CategoryScores       JSONB     `json:"categoryScores" validate:"required"`
	WeightedAverageScore float64   `json:"weightedAverageScore" validate:"required,min=0,max=10"`
	Highlights           string    `json:"highlights"`
	PointsToDevelop      string    `json:"pointsToDevelop"`
}

type ArchiveDeveloperRequest struct {
	Archive bool `json:"archive"`
}

type User struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	Email               string    `json:"email" db:"email"`
	Password            string    `json:"-" db:"password"` // O "-" faz com que este campo não seja serializado no JSON
	Name                string    `json:"name" db:"name"`
	Role                string    `json:"role" db:"role"` // admin, manager, user
	NeedsPasswordChange bool      `json:"needsPasswordChange" db:"needs_password_change"`
	IsActive            bool      `json:"isActive" db:"is_active"`
	CreatedAt           time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time `json:"updatedAt" db:"updated_at"`
}

func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required,min=2"`
	Role     string `json:"role" validate:"omitempty,oneof=admin manager user"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type JWTClaims struct {
	UserID              uuid.UUID `json:"userId"`
	Email               string    `json:"email"`
	Role                string    `json:"role"`
	IsActive            bool      `json:"isActive"`
	NeedsPasswordChange bool      `json:"needsPasswordChange"`
}

type CreateUserRequest struct {
	Name              string `json:"name" validate:"required,min=2"`
	Email             string `json:"email" validate:"required,email"`
	Role              string `json:"role" validate:"required,oneof=admin manager user"`
	TemporaryPassword string `json:"temporaryPassword" validate:"required,min=8"`
}

type SetNewPasswordRequest struct {
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=8"`
}
