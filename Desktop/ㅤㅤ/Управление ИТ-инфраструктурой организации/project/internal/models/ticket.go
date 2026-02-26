package models

import "time"

type Ticket struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Category    string    `json:"category" db:"category"`
	Location    string    `json:"location" db:"location"`
	Description string    `json:"description" db:"description"`
	PhotoURL    string    `json:"photo_url,omitempty" db:"photo_url"`
	Status      string    `json:"status" db:"status"`
	UserID      int       `json:"user_id" db:"user_id"`
	AssigneeID  *int      `json:"assignee_id,omitempty" db:"assignee_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateTicketRequest struct {
	Title       string `json:"title" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Location    string `json:"location" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=new in_progress done"`
}
