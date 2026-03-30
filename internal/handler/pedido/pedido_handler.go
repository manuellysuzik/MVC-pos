package pedidohdl

import (
	"errors"
	"net/http"
	"strconv"

	"ecommerce-api/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

type pedidoService interface {
	Create(p *entity.Pedido) error
	FindAll() ([]entity.Pedido, error)
	FindByID(id uint) (*entity.Pedido, error)
	FindByStatus(status string) ([]entity.Pedido, error)
	FindByClienteID(clienteID uint) ([]entity.Pedido, error)
	Count() (int64, error)
	Update(p *entity.Pedido) error
	Delete(id uint) error
}

type PedidoHandler struct{ service pedidoService }

func New(svc pedidoService) *PedidoHandler {
	return &PedidoHandler{service: svc}
}

func handleServiceError(c *gin.Context, err error) {
	if errors.Is(err, entity.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var valErr *entity.ValidationError
	if errors.As(err, &valErr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": valErr.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}

func (h *PedidoHandler) Create(c *gin.Context) {
	var p entity.Pedido
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(&p); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *PedidoHandler) FindAll(c *gin.Context) {
	pedidos, err := h.service.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, pedidos)
}

func (h *PedidoHandler) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
	p, err := h.service.FindByID(uint(id))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *PedidoHandler) FindByStatus(c *gin.Context) {
	status := c.Query("status")
	pedidos, err := h.service.FindByStatus(status)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, pedidos)
}

func (h *PedidoHandler) FindByClienteID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
	pedidos, err := h.service.FindByClienteID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, pedidos)
}

func (h *PedidoHandler) Count(c *gin.Context) {
	count, err := h.service.Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (h *PedidoHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
	var p entity.Pedido
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.ID = uint(id)
	if err := h.service.Update(&p); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *PedidoHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
	if err := h.service.Delete(uint(id)); err != nil {
		handleServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
