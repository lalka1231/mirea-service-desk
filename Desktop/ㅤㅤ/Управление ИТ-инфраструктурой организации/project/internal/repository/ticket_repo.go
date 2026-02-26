package repository

import (
	"database/sql"
	"mirea-service-desk/internal/models"
)

type TicketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) Create(ticket *models.Ticket) error {
	query := `
		INSERT INTO tickets (title, category, location, description, photo_url, status, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query,
		ticket.Title, ticket.Category, ticket.Location, ticket.Description,
		ticket.PhotoURL, "new", ticket.UserID,
	).Scan(&ticket.ID, &ticket.CreatedAt)
}

func (r *TicketRepository) GetAll(userID int, role string) ([]models.Ticket, error) {
	var rows *sql.Rows
	var err error
	
	if role == "executor" {
		rows, err = r.db.Query(`
			SELECT id, title, category, location, description, photo_url, status, user_id, assignee_id, created_at, updated_at
			FROM tickets ORDER BY created_at DESC
		`)
	} else {
		rows, err = r.db.Query(`
			SELECT id, title, category, location, description, photo_url, status, user_id, assignee_id, created_at, updated_at
			FROM tickets WHERE user_id = $1 ORDER BY created_at DESC
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
		FROM tickets WHERE id = $1
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
	query := `UPDATE tickets SET status = $1, assignee_id = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.Exec(query, status, assigneeID, id)
	return err
}
