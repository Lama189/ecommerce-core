package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var uzbPhoneRegex = regexp.MustCompile(`^\+998\d{9}$`)

type Role string

const (
	RoleUser   Role = "user"
	RoleArtist Role = "artist"
	RoleAdmin  Role = "admin"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleUser, RoleArtist, RoleAdmin:
		return true
	default:
		return false
	}
}

type User struct {
	ID           uuid.UUID
	Phone        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(phone, passwordHash string) (*User, error) {
	now := time.Now().UTC()
	u := &User{
		Phone:        strings.TrimSpace(phone),
		PasswordHash: passwordHash,
		Role:         RoleUser,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := u.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}

	return u, nil
}

func RestoreUser(id uuid.UUID, phone, passwordHash string, role Role, createdAt, updatedAt time.Time) *User {
	if role == "" {
		role = RoleUser
	}

	return &User{
		ID:           id,
		Phone:        phone,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

func (u *User) Validate() error {
	if u.Phone == "" {
		return fmt.Errorf("phone number cannot be emplty")
	}

	if len(u.Phone) != 13 {
		return fmt.Errorf("invalid phone length: expected 13 characters, got %d", len(u.Phone))
	}

	if !strings.HasPrefix(u.Phone, "+998") {
		return fmt.Errorf("invalid Uzbekistan phone format")
	}

	if len(u.PasswordHash) == 0 {
		return fmt.Errorf("password hash cannot be empty")
	}

	if !u.Role.IsValid() {
		return fmt.Errorf("invalid user role: %s", u.Role)
	}

	return nil
}
