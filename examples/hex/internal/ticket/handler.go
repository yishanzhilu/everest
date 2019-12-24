package ticket

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler ...
type Handler interface {
	Get(c *gin.Context)
	GetByID(c *gin.Context)
	Create(c *gin.Context)
}

type ticketHandler struct {
	ticketService Service
}

// NewTicketHandler ...
func NewTicketHandler(ticketService Service) Handler {
	return &ticketHandler{
		ticketService,
	}
}

func (h *ticketHandler) Get(c *gin.Context) {
	tickets, _ := h.ticketService.FindAllTickets()
	c.JSON(http.StatusOK, tickets)
}

func (h *ticketHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	ticket, _ := h.ticketService.FindTicketByID(id)
	c.JSON(http.StatusOK, ticket)
}

func (h *ticketHandler) Create(c *gin.Context) {

	var ticket Model
	if err := c.Bind(&ticket); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	_ = h.ticketService.CreateTicket(&ticket)
	c.JSON(http.StatusOK, ticket)
}
