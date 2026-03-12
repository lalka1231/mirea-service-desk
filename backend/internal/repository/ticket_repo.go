package repository

import (
	"database/sql"

	"github.com/lalka1231/mirea-service-desk/internal/models"
)

type TicketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) Create(ticket *models.Ticket) error {
	// SQLite не поддерживает RETURNING, используем LastInsertId
	query := `
		INSERT INTO tickets (title, category, location, description, photo_url, status, user_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	result, err := r.db.Exec(query,
		ticket.Title, ticket.Category, ticket.Location, ticket.Description,
		ticket.PhotoURL, "new", ticket.UserID,
	)
	if err != nil {
		return err
	}

	// Получаем ID новой записи
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	ticket.ID = int(id)

	// Получаем created_at
	err = r.db.QueryRow("SELECT created_at FROM tickets WHERE id = ?", ticket.ID).Scan(&ticket.CreatedAt)
	return err
}

func (r *TicketRepository) GetAll(userID int, role string) ([]models.Ticket, error) {
	var rows *sql.Rows
	var err error

	if role == "executor" {
		// Исполнитель видит все заявки
		rows, err = r.db.Query(`
			SELECT id, title, category, location, description, photo_url, status, user_id, assignee_id, created_at, updated_at
			FROM tickets ORDER BY created_at DESC
		`)
	} else {
		// Студент видит только свои
		rows, err = r.db.Query(`
			SELECT id, title, category, location, description, photo_url, status, user_id, assignee_id, created_at, updated_at
			FROM tickets WHERE user_id = ? ORDER BY created_at DESC
		`, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []models.Ticket
	for rows.Next() {
		var t models.Ticket
		err := rows.Scan(
			&t.ID, &t.Title, &t.Category, &t.Location, &t.Description,
			&t.PhotoURL, &t.Status, &t.UserID, &t.AssigneeID, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	return tickets, nil
}

func (r *TicketRepository) GetByID(id int) (*models.Ticket, error) {
	var t models.Ticket
	query := `
		SELECT id, title, category, location, description, photo_url, status, user_id, assignee_id, created_at, updated_at
		FROM tickets WHERE id = ?
	`
	err := r.db.QueryRow(query, id).Scan(
		&t.ID, &t.Title, &t.Category, &t.Location, &t.Description,
		&t.PhotoURL, &t.Status, &t.UserID, &t.AssigneeID, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TicketRepository) UpdateStatus(id int, status string, assigneeID *int) error {
	query := `UPDATE tickets SET status = ?, assignee_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, status, assigneeID, id)
	return err
}
