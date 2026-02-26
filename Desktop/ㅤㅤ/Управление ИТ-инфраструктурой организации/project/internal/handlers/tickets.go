package handlers

import (
	"net/http"
	"strconv"
	"github.com/lalka1231/mirea-service-desk/internal/models"
	"github.com/lalka1231/mirea-service-desk/internal/repository"
	"github.com/lalka1231/mirea-service-desk/pkg/utils"
	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	ticketRepo *repository.TicketRepository
}

func NewTicketHandler(ticketRepo *repository.TicketRepository) *TicketHandler {
	return &TicketHandler{ticketRepo: ticketRepo}
}

func (h *TicketHandler) Create(c *gin.Context) {
	userID, _ := utils.GetUserFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	var req models.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	ticket := &models.Ticket{
		Title:       req.Title,
		Category:    req.Category,
		Location:    req.Location,
		Description: req.Description,
		UserID:      userID,
	}
	
	if err := h.ticketRepo.Create(ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
		return
	}
	
	c.JSON(http.StatusCreated, ticket)
}

func (h *TicketHandler) GetAll(c *gin.Context) {
	userID, role := utils.GetUserFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	tickets, err := h.ticketRepo.GetAll(userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tickets"})
		return
	}
	
	c.JSON(http.StatusOK, tickets)
}

func (h *TicketHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}
	
	ticket, err := h.ticketRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}
	
	c.JSON(http.StatusOK, ticket)
}

func (h *TicketHandler) UpdateStatus(c *gin.Context) {
	userID, role := utils.GetUserFromContext(c)
	if role != "executor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only executors can update status"})
		return
	}
	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}
	
	var req models.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	var assigneeID *int
	if req.Status == "in_progress" {
		assigneeID = &userID
	}
	
	if err := h.ticketRepo.UpdateStatus(id, req.Status, assigneeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}
