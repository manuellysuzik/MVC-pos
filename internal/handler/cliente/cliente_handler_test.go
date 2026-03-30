package clienthdl

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

type mockClienteSvc struct {
	clientes map[uint]*entity.Cliente
	nextID   uint
}

func newMockClienteSvc() *mockClienteSvc {
	return &mockClienteSvc{clientes: make(map[uint]*entity.Cliente), nextID: 1}
}
func (m *mockClienteSvc) Create(c *entity.Cliente) error {
	if c.Nome == "" {
		return entity.NewValidationError("nome é obrigatório")
	}
	c.ID = m.nextID
	m.nextID++
	m.clientes[c.ID] = c
	return nil
}
func (m *mockClienteSvc) FindAll() ([]entity.Cliente, error) {
	r := make([]entity.Cliente, 0)
	for _, c := range m.clientes {
		r = append(r, *c)
	}
	return r, nil
}
func (m *mockClienteSvc) FindByID(id uint) (*entity.Cliente, error) {
	c, ok := m.clientes[id]
	if !ok {
		return nil, entity.ErrNotFound
	}
	return c, nil
}
func (m *mockClienteSvc) FindByName(nome string) ([]entity.Cliente, error) {
	if nome == "" {
		return nil, errors.New("nome é obrigatório para busca")
	}
	return []entity.Cliente{}, nil
}
func (m *mockClienteSvc) Count() (int64, error) { return int64(len(m.clientes)), nil }
func (m *mockClienteSvc) Update(c *entity.Cliente) error {
	if _, ok := m.clientes[c.ID]; !ok {
		return entity.ErrNotFound
	}
	m.clientes[c.ID] = c
	return nil
}
func (m *mockClienteSvc) Delete(id uint) error {
	if _, ok := m.clientes[id]; !ok {
		return entity.ErrNotFound
	}
	delete(m.clientes, id)
	return nil
}

func setupClienteRouter(svc clienteService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := New(svc)
	r.POST("/clientes", h.Create)
	r.GET("/clientes", h.FindAll)
	r.GET("/clientes/count", h.Count)
	r.GET("/clientes/search", h.FindByName)
	r.GET("/clientes/:id", h.FindByID)
	r.PUT("/clientes/:id", h.Update)
	r.DELETE("/clientes/:id", h.Delete)
	return r
}

func TestClienteHandler_Create(t *testing.T) {
	r := setupClienteRouter(newMockClienteSvc())
	body, _ := json.Marshal(map[string]string{"nome": "João", "email": "joao@test.com"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/clientes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestClienteHandler_Create_BadBody(t *testing.T) {
	r := setupClienteRouter(newMockClienteSvc())
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/clientes", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestClienteHandler_FindAll(t *testing.T) {
	svc := newMockClienteSvc()
	svc.Create(&entity.Cliente{Nome: "Ana", Email: "ana@test.com"})
	r := setupClienteRouter(svc)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/clientes", nil))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestClienteHandler_FindByID(t *testing.T) {
	svc := newMockClienteSvc()
	svc.Create(&entity.Cliente{Nome: "Bob", Email: "bob@test.com"})
	r := setupClienteRouter(svc)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/clientes/1", nil))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestClienteHandler_FindByID_NotFound(t *testing.T) {
	r := setupClienteRouter(newMockClienteSvc())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/clientes/999", nil))
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestClienteHandler_Count(t *testing.T) {
	r := setupClienteRouter(newMockClienteSvc())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/clientes/count", nil))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestClienteHandler_Delete_NotFound(t *testing.T) {
	r := setupClienteRouter(newMockClienteSvc())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("DELETE", "/clientes/999", nil))
	assert.Equal(t, http.StatusNotFound, w.Code)
}
