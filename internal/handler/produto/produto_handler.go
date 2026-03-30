package produtohdl

import (
	"errors"
	"net/http"
	"strconv"

	"ecommerce-api/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

type produtoService interface {
	Create(p *entity.Produto) error
	FindAll() ([]entity.Produto, error)
	FindByID(id uint) (*entity.Produto, error)
	FindByName(nome string) ([]entity.Produto, error)
	Count() (int64, error)
	Update(p *entity.Produto) error
	Delete(id uint) error
}

type ProdutoHandler struct{ service produtoService }

func New(svc produtoService) *ProdutoHandler {
	return &ProdutoHandler{service: svc}
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

func (h *ProdutoHandler) Create(c *gin.Context) {
	var p entity.Produto
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

func (h *ProdutoHandler) FindAll(c *gin.Context) {
	produtos, err := h.service.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, produtos)
}

func (h *ProdutoHandler) FindByID(c *gin.Context) {
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

func (h *ProdutoHandler) FindByName(c *gin.Context) {
	nome := c.Query("nome")
	produtos, err := h.service.FindByName(nome)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, produtos)
}

func (h *ProdutoHandler) Count(c *gin.Context) {
	count, err := h.service.Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (h *ProdutoHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
	var p entity.Produto
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

func (h *ProdutoHandler) Delete(c *gin.Context) {
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
