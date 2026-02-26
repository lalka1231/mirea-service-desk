package repository

import (
	"database/sql"
	"mirea-service-desk/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, role, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, created_at
	`
	return r.db.QueryRow(query, user.Email, user.Password, user.Role).Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, password_hash, role, created_at FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, role, created_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
