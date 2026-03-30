package clienthdl

import (
	"errors"
	"net/http"
	"strconv"

	"ecommerce-api/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

type clienteService interface {
	Create(c *entity.Cliente) error
	FindAll() ([]entity.Cliente, error)
	FindByID(id uint) (*entity.Cliente, error)
	FindByName(nome string) ([]entity.Cliente, error)
	Count() (int64, error)
	Update(c *entity.Cliente) error
	Delete(id uint) error
}

type ClienteHandler struct{ service clienteService }

func New(svc clienteService) *ClienteHandler {
	return &ClienteHandler{service: svc}
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

func (h *ClienteHandler) Create(c *gin.Context) {
	var cliente entity.Cliente
	if err := c.ShouldBindJSON(&cliente); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(&cliente); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, cliente)
}

func (h *ClienteHandler) FindAll(c *gin.Context) {
	clientes, err := h.service.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, clientes)
}

func (h *ClienteHandler) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
	cliente, err := h.service.FindByID(uint(id))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, cliente)
}

func (h *ClienteHandler) FindByName(c *gin.Context) {
	nome := c.Query("nome")
	clientes, err := h.service.FindByName(nome)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, clientes)
}

func (h *ClienteHandler) Count(c *gin.Context) {
	count, err := h.service.Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (h *ClienteHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
	var cliente entity.Cliente
	if err := c.ShouldBindJSON(&cliente); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cliente.ID = uint(id)
	if err := h.service.Update(&cliente); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, cliente)
}

func (h *ClienteHandler) Delete(c *gin.Context) {
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
