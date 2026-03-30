package produtohdl

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"ecommerce-api/internal/domain/entity"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockProdutoSvc struct {
	produtos map[uint]*entity.Produto
	nextID   uint
}

func newMockProdutoSvc() *mockProdutoSvc {
	return &mockProdutoSvc{produtos: make(map[uint]*entity.Produto), nextID: 1}
}
func (m *mockProdutoSvc) Create(p *entity.Produto) error {
	if p.Nome == "" {
		return entity.NewValidationError("nome é obrigatório")
	}
	p.ID = m.nextID
	m.nextID++
	m.produtos[p.ID] = p
	return nil
}
func (m *mockProdutoSvc) FindAll() ([]entity.Produto, error) {
	r := make([]entity.Produto, 0)
	for _, p := range m.produtos {
		r = append(r, *p)
	}
	return r, nil
}
func (m *mockProdutoSvc) FindByID(id uint) (*entity.Produto, error) {
	p, ok := m.produtos[id]
	if !ok {
		return nil, entity.ErrNotFound
	}
	return p, nil
}
func (m *mockProdutoSvc) FindByName(nome string) ([]entity.Produto, error) {
	if nome == "" {
		return nil, errors.New("nome é obrigatório para busca")
	}
	return []entity.Produto{}, nil
}
func (m *mockProdutoSvc) Count() (int64, error) { return int64(len(m.produtos)), nil }
func (m *mockProdutoSvc) Update(p *entity.Produto) error {
	if _, ok := m.produtos[p.ID]; !ok {
		return entity.ErrNotFound
	}
	m.produtos[p.ID] = p
	return nil
}
func (m *mockProdutoSvc) Delete(id uint) error {
	if _, ok := m.produtos[id]; !ok {
		return entity.ErrNotFound
	}
	delete(m.produtos, id)
	return nil
}

func setupProdutoRouter(svc produtoService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := New(svc)
	r.POST("/produtos", h.Create)
	r.GET("/produtos", h.FindAll)
	r.GET("/produtos/count", h.Count)
	r.GET("/produtos/search", h.FindByName)
	r.GET("/produtos/:id", h.FindByID)
	r.PUT("/produtos/:id", h.Update)
	r.DELETE("/produtos/:id", h.Delete)
	return r
}

func TestProdutoHandler_Create(t *testing.T) {
	r := setupProdutoRouter(newMockProdutoSvc())
	body, _ := json.Marshal(map[string]interface{}{"nome": "Notebook", "preco": 3500.00})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/produtos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestProdutoHandler_FindAll(t *testing.T) {
	svc := newMockProdutoSvc()
	svc.Create(&entity.Produto{Nome: "Mouse", Preco: 80.00})
	r := setupProdutoRouter(svc)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/produtos", nil))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProdutoHandler_FindByID_NotFound(t *testing.T) {
	r := setupProdutoRouter(newMockProdutoSvc())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/produtos/999", nil))
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestProdutoHandler_Count(t *testing.T) {
	r := setupProdutoRouter(newMockProdutoSvc())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/produtos/count", nil))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProdutoHandler_Delete_NotFound(t *testing.T) {
	r := setupProdutoRouter(newMockProdutoSvc())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("DELETE", "/produtos/999", nil))
	assert.Equal(t, http.StatusNotFound, w.Code)
}
