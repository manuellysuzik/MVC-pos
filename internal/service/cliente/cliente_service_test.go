package clientesvc

import (
	"testing"

	"ecommerce-api/internal/domain/entity"

	"github.com/stretchr/testify/assert"
)

type mockClienteRepo struct {
	clientes map[uint]*entity.Cliente
	nextID   uint
}

func newMockClienteRepo() *mockClienteRepo {
	return &mockClienteRepo{clientes: make(map[uint]*entity.Cliente), nextID: 1}
}
func (m *mockClienteRepo) Create(c *entity.Cliente) error {
	c.ID = m.nextID
	m.nextID++
	m.clientes[c.ID] = c
	return nil
}
func (m *mockClienteRepo) FindAll() ([]entity.Cliente, error) {
	r := make([]entity.Cliente, 0)
	for _, c := range m.clientes {
		r = append(r, *c)
	}
	return r, nil
}
func (m *mockClienteRepo) FindByID(id uint) (*entity.Cliente, error) {
	c, ok := m.clientes[id]
	if !ok {
		return nil, entity.ErrNotFound
	}
	return c, nil
}
func (m *mockClienteRepo) FindByName(nome string) ([]entity.Cliente, error) {
	return []entity.Cliente{}, nil
}
func (m *mockClienteRepo) Count() (int64, error)          { return int64(len(m.clientes)), nil }
func (m *mockClienteRepo) Update(c *entity.Cliente) error { m.clientes[c.ID] = c; return nil }
func (m *mockClienteRepo) Delete(id uint) error           { delete(m.clientes, id); return nil }

func TestClienteSvc_Create_Success(t *testing.T) {
	svc := New(newMockClienteRepo())
	c := &entity.Cliente{Nome: "João", Email: "joao@test.com"}
	assert.NoError(t, svc.Create(c))
	assert.NotZero(t, c.ID)
}

func TestClienteSvc_Create_MissingNome(t *testing.T) {
	svc := New(newMockClienteRepo())
	var valErr *entity.ValidationError
	assert.ErrorAs(t, svc.Create(&entity.Cliente{Email: "x@test.com"}), &valErr)
	assert.Equal(t, "nome é obrigatório", valErr.Message)
}

func TestClienteSvc_Create_MissingEmail(t *testing.T) {
	svc := New(newMockClienteRepo())
	var valErr *entity.ValidationError
	assert.ErrorAs(t, svc.Create(&entity.Cliente{Nome: "X"}), &valErr)
	assert.Equal(t, "email é obrigatório", valErr.Message)
}

func TestClienteSvc_FindAll(t *testing.T) {
	repo := newMockClienteRepo()
	svc := New(repo)
	svc.Create(&entity.Cliente{Nome: "A", Email: "a@test.com"})
	svc.Create(&entity.Cliente{Nome: "B", Email: "b@test.com"})
	clientes, err := svc.FindAll()
	assert.NoError(t, err)
	assert.Len(t, clientes, 2)
}

func TestClienteSvc_FindByName_EmptyName(t *testing.T) {
	svc := New(newMockClienteRepo())
	_, err := svc.FindByName("")
	var valErr *entity.ValidationError
	assert.ErrorAs(t, err, &valErr)
}

func TestClienteSvc_Update_NotFound(t *testing.T) {
	svc := New(newMockClienteRepo())
	err := svc.Update(&entity.Cliente{ID: 99, Nome: "X", Email: "x@test.com"})
	assert.ErrorIs(t, err, entity.ErrNotFound)
}

func TestClienteSvc_Delete_NotFound(t *testing.T) {
	svc := New(newMockClienteRepo())
	assert.ErrorIs(t, svc.Delete(99), entity.ErrNotFound)
}
