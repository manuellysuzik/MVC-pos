package produtosvc

import (
	"testing"

	"ecommerce-api/internal/domain/entity"

	"github.com/stretchr/testify/assert"
)

type mockProdutoRepo struct {
	produtos map[uint]*entity.Produto
	nextID   uint
}

func newMockProdutoRepo() *mockProdutoRepo {
	return &mockProdutoRepo{produtos: make(map[uint]*entity.Produto), nextID: 1}
}
func (m *mockProdutoRepo) Create(p *entity.Produto) error {
	p.ID = m.nextID
	m.nextID++
	m.produtos[p.ID] = p
	return nil
}
func (m *mockProdutoRepo) FindAll() ([]entity.Produto, error) {
	r := make([]entity.Produto, 0)
	for _, p := range m.produtos {
		r = append(r, *p)
	}
	return r, nil
}
func (m *mockProdutoRepo) FindByID(id uint) (*entity.Produto, error) {
	p, ok := m.produtos[id]
	if !ok {
		return nil, entity.ErrNotFound
	}
	return p, nil
}
func (m *mockProdutoRepo) FindByName(nome string) ([]entity.Produto, error) {
	return []entity.Produto{}, nil
}
func (m *mockProdutoRepo) Count() (int64, error)          { return int64(len(m.produtos)), nil }
func (m *mockProdutoRepo) Update(p *entity.Produto) error { m.produtos[p.ID] = p; return nil }
func (m *mockProdutoRepo) Delete(id uint) error           { delete(m.produtos, id); return nil }

func TestProdutoSvc_Create_Success(t *testing.T) {
	svc := New(newMockProdutoRepo())
	p := &entity.Produto{Nome: "Notebook", Preco: 3500.00}
	assert.NoError(t, svc.Create(p))
	assert.NotZero(t, p.ID)
}

func TestProdutoSvc_Create_MissingNome(t *testing.T) {
	svc := New(newMockProdutoRepo())
	var valErr *entity.ValidationError
	assert.ErrorAs(t, svc.Create(&entity.Produto{Preco: 10.00}), &valErr)
	assert.Equal(t, "nome é obrigatório", valErr.Message)
}

func TestProdutoSvc_Create_PrecoInvalido(t *testing.T) {
	svc := New(newMockProdutoRepo())
	var valErr *entity.ValidationError
	assert.ErrorAs(t, svc.Create(&entity.Produto{Nome: "X", Preco: -1}), &valErr)
	assert.Equal(t, "preço deve ser maior que zero", valErr.Message)
}

func TestProdutoSvc_FindByName_EmptyName(t *testing.T) {
	svc := New(newMockProdutoRepo())
	_, err := svc.FindByName("")
	var valErr *entity.ValidationError
	assert.ErrorAs(t, err, &valErr)
}

func TestProdutoSvc_Update_NotFound(t *testing.T) {
	svc := New(newMockProdutoRepo())
	err := svc.Update(&entity.Produto{ID: 99, Nome: "X", Preco: 1.00})
	assert.ErrorIs(t, err, entity.ErrNotFound)
}

func TestProdutoSvc_Delete_NotFound(t *testing.T) {
	svc := New(newMockProdutoRepo())
	assert.ErrorIs(t, svc.Delete(99), entity.ErrNotFound)
}
