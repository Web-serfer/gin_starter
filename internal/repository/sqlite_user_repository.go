package repository

import (
	"database/sql"
	"fmt"
	"gin-starter/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteUserRepository реализация репозитория для SQLite
type SQLiteUserRepository struct {
	db *sql.DB
}

// NewSQLiteUserRepository создает новый экземпляр репозитория
func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{
		db: db,
	}
}

// Create создает нового пользователя
func (r *SQLiteUserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (name, email, created_at, updated_at)
		VALUES (?, ?, datetime('now'), datetime('now'))
	`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	result, err := stmt.Exec(user.Name, user.Email)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = uint(id)

	// Устанавливаем даты создания и обновления
	row := r.db.QueryRow("SELECT created_at, updated_at FROM users WHERE id = ?", user.ID)
	err = row.Scan(&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to get user dates: %w", err)
	}

	return nil
}

// GetByID возвращает пользователя по ID
func (r *SQLiteUserRepository) GetByID(id uint) (*models.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?`

	row := r.db.QueryRow(query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByEmail возвращает пользователя по email
func (r *SQLiteUserRepository) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE email = ?`

	row := r.db.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetAll возвращает всех пользователей
func (r *SQLiteUserRepository) GetAll() ([]*models.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	return users, nil
}

// Update обновляет пользователя
func (r *SQLiteUserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET name = ?, email = ?, updated_at = datetime('now')
		WHERE id = ?
	`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(user.Name, user.Email, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete удаляет пользователя по ID
func (r *SQLiteUserRepository) Delete(id uint) error {
	query := `DELETE FROM users WHERE id = ?`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
