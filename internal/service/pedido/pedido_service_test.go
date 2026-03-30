package pedidosvc

import (
	"testing"

	"ecommerce-api/internal/domain/entity"

	"github.com/stretchr/testify/assert"
)

type mockPedidoRepo struct {
	pedidos map[uint]*entity.Pedido
	nextID  uint
}

func newMockPedidoRepo() *mockPedidoRepo {
	return &mockPedidoRepo{pedidos: make(map[uint]*entity.Pedido), nextID: 1}
}
func (m *mockPedidoRepo) Create(p *entity.Pedido) error {
	p.ID = m.nextID
	m.nextID++
	m.pedidos[p.ID] = p
	return nil
}
func (m *mockPedidoRepo) FindAll() ([]entity.Pedido, error) {
	r := make([]entity.Pedido, 0)
	for _, p := range m.pedidos {
		r = append(r, *p)
	}
	return r, nil
}
func (m *mockPedidoRepo) FindByID(id uint) (*entity.Pedido, error) {
	p, ok := m.pedidos[id]
	if !ok {
		return nil, entity.ErrNotFound
	}
	return p, nil
}
func (m *mockPedidoRepo) FindByStatus(status string) ([]entity.Pedido, error) {
	return []entity.Pedido{}, nil
}
func (m *mockPedidoRepo) FindByClienteID(clienteID uint) ([]entity.Pedido, error) {
	return []entity.Pedido{}, nil
}
func (m *mockPedidoRepo) Count() (int64, error)         { return int64(len(m.pedidos)), nil }
func (m *mockPedidoRepo) Update(p *entity.Pedido) error { m.pedidos[p.ID] = p; return nil }
func (m *mockPedidoRepo) Delete(id uint) error          { delete(m.pedidos, id); return nil }

func TestPedidoSvc_Create_Success(t *testing.T) {
	svc := New(newMockPedidoRepo())
	p := &entity.Pedido{
		ClienteID: 1,
		Itens:     []entity.ItemPedido{{ProdutoID: 1, Quantidade: 2, PrecoUnitario: 50.00}},
	}
	assert.NoError(t, svc.Create(p))
	assert.Equal(t, 100.00, p.Total)
	assert.Equal(t, "pendente", p.Status)
}

func TestPedidoSvc_Create_MissingClienteID(t *testing.T) {
	svc := New(newMockPedidoRepo())
	err := svc.Create(&entity.Pedido{Itens: []entity.ItemPedido{{ProdutoID: 1, Quantidade: 1, PrecoUnitario: 10.00}}})
	var valErr *entity.ValidationError
	assert.ErrorAs(t, err, &valErr)
}

func TestPedidoSvc_Create_SemItens(t *testing.T) {
	svc := New(newMockPedidoRepo())
	err := svc.Create(&entity.Pedido{ClienteID: 1})
	var valErr *entity.ValidationError
	assert.ErrorAs(t, err, &valErr)
}

func TestPedidoSvc_Create_QuantidadeInvalida(t *testing.T) {
	svc := New(newMockPedidoRepo())
	err := svc.Create(&entity.Pedido{
		ClienteID: 1,
		Itens:     []entity.ItemPedido{{ProdutoID: 1, Quantidade: 0, PrecoUnitario: 10.00}},
	})
	var valErr *entity.ValidationError
	assert.ErrorAs(t, err, &valErr)
}

func TestPedidoSvc_Update_NotFound(t *testing.T) {
	svc := New(newMockPedidoRepo())
	err := svc.Update(&entity.Pedido{ID: 99, ClienteID: 1, Status: "confirmado"})
	assert.ErrorIs(t, err, entity.ErrNotFound)
}

func TestPedidoSvc_Delete_NotFound(t *testing.T) {
	svc := New(newMockPedidoRepo())
	assert.ErrorIs(t, svc.Delete(99), entity.ErrNotFound)
}
