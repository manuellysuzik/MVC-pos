package pedidohdl

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

type mockPedidoSvc struct {
	pedidos map[uint]*entity.Pedido
	nextID  uint
}

func newMockPedidoSvc() *mockPedidoSvc {
	return &mockPedidoSvc{pedidos: make(map[uint]*entity.Pedido), nextID: 1}
}
func (m *mockPedidoSvc) Create(p *entity.Pedido) error {
	if p.ClienteID == 0 {
		return entity.NewValidationError("cliente_id é obrigatório")
	}
	p.ID = m.nextID
	m.nextID++
	m.pedidos[p.ID] = p
	return nil
}
func (m *mockPedidoSvc) FindAll() ([]entity.Pedido, error) {
	r := make([]entity.Pedido, 0)
	for _, p := range m.pedidos {
		r = append(r, *p)
	}
	return r, nil
}
func (m *mockPedidoSvc) FindByID(id uint) (*entity.Pedido, error) {
	p, ok := m.pedidos[id]
	if !ok {
		return nil, entity.ErrNotFound
	}
	return p, nil
}
func (m *mockPedidoSvc) FindByStatus(status string) ([]entity.Pedido, error) {
	if status == "" {
		return nil, errors.New("status é obrigatório para busca")
	}
	return []entity.Pedido{}, nil
}
func (m *mockPedidoSvc) FindByClienteID(clienteID uint) ([]entity.Pedido, error) {
	return []entity.Pedido{}, nil
}
func (m *mockPedidoSvc) Count() (int64, error) { return int64(len(m.pedidos)), nil }
func (m *mockPedidoSvc) Update(p *entity.Pedido) error {
	if _, ok := m.pedidos[p.ID]; !ok {
		return entity.ErrNotFound
	}
	m.pedidos[p.ID] = p
	return nil
}
func (m *mockPedidoSvc) Delete(id uint) error {
	if _, ok := m.pedidos[id]; !ok {
		return entity.ErrNotFound
	}
	delete(m.pedidos, id)
	return nil
}

func setupPedidoRouter(svc pedidoService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := New(svc)
	r.POST("/pedidos", h.Create)
	r.GET("/pedidos", h.FindAll)
	r.GET("/pedidos/count", h.Count)
	r.GET("/pedidos/search", h.FindByStatus)
	r.GET("/pedidos/cliente/:id", h.FindByClienteID)
	r.GET("/pedidos/:id", h.FindByID)
	r.PUT("/pedidos/:id", h.Update)
	r.DELETE("/pedidos/:id", h.Delete)
	return r
}

func TestPedidoHandler_Create(t *testing.T) {
	r := setupPedidoRouter(newMockPedidoSvc())
	body, _ := json.Marshal(map[string]interface{}{
		"cliente_id": 1,
		"itens":      []map[string]interface{}{{"produto_id": 1, "quantidade": 2, "preco_unitario": 50.00}},
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/pedidos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPedidoHandler_FindAll(t *testing.T) {
	r := setupPedidoRouter(newMockPedidoSvc())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/pedidos", nil))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPedidoHandler_FindByID_NotFound(t *testing.T) {
	r := setupPedidoRouter(newMockPedidoSvc())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/pedidos/999", nil))
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPedidoHandler_Count(t *testing.T) {
	r := setupPedidoRouter(newMockPedidoSvc())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/pedidos/count", nil))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPedidoHandler_FindByClienteID(t *testing.T) {
	r := setupPedidoRouter(newMockPedidoSvc())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/pedidos/cliente/1", nil))
	assert.Equal(t, http.StatusOK, w.Code)
}
