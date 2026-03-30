package pedidorepo

import (
	"testing"

	"ecommerce-api/internal/domain/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupPedidoDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&entity.Cliente{}, &entity.Produto{}, &entity.Pedido{}, &entity.ItemPedido{}))
	return db
}

func seedClienteAndProduto(t *testing.T, db *gorm.DB) (entity.Cliente, entity.Produto) {
	c := entity.Cliente{Nome: "Cliente Teste", Email: "ct@test.com"}
	require.NoError(t, db.Create(&c).Error)
	p := entity.Produto{Nome: "Produto Teste", Preco: 100.00}
	require.NoError(t, db.Create(&p).Error)
	return c, p
}

func TestPedidoRepo_Create(t *testing.T) {
	db := setupPedidoDB(t)
	c, p := seedClienteAndProduto(t, db)
	repo := New(db)
	pedido := &entity.Pedido{
		ClienteID: c.ID,
		Status:    "pendente",
		Total:     100.00,
		Itens:     []entity.ItemPedido{{ProdutoID: p.ID, Quantidade: 1, PrecoUnitario: 100.00}},
	}
	assert.NoError(t, repo.Create(pedido))
	assert.NotZero(t, pedido.ID)
}

func TestPedidoRepo_FindAll(t *testing.T) {
	db := setupPedidoDB(t)
	c, _ := seedClienteAndProduto(t, db)
	repo := New(db)
	repo.Create(&entity.Pedido{ClienteID: c.ID, Status: "pendente", Total: 50.00})
	repo.Create(&entity.Pedido{ClienteID: c.ID, Status: "confirmado", Total: 200.00})
	pedidos, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, pedidos, 2)
}

func TestPedidoRepo_FindByID(t *testing.T) {
	db := setupPedidoDB(t)
	c, _ := seedClienteAndProduto(t, db)
	repo := New(db)
	ped := &entity.Pedido{ClienteID: c.ID, Status: "pendente", Total: 100.00}
	repo.Create(ped)
	found, err := repo.FindByID(ped.ID)
	assert.NoError(t, err)
	assert.Equal(t, ped.ID, found.ID)
}

func TestPedidoRepo_FindByStatus(t *testing.T) {
	db := setupPedidoDB(t)
	c, _ := seedClienteAndProduto(t, db)
	repo := New(db)
	repo.Create(&entity.Pedido{ClienteID: c.ID, Status: "pendente", Total: 50.00})
	repo.Create(&entity.Pedido{ClienteID: c.ID, Status: "confirmado", Total: 200.00})
	repo.Create(&entity.Pedido{ClienteID: c.ID, Status: "pendente", Total: 75.00})
	results, err := repo.FindByStatus("pendente")
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestPedidoRepo_FindByClienteID(t *testing.T) {
	db := setupPedidoDB(t)
	c, _ := seedClienteAndProduto(t, db)
	repo := New(db)
	repo.Create(&entity.Pedido{ClienteID: c.ID, Status: "pendente", Total: 50.00})
	repo.Create(&entity.Pedido{ClienteID: c.ID, Status: "confirmado", Total: 200.00})
	results, err := repo.FindByClienteID(c.ID)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestPedidoRepo_Count(t *testing.T) {
	db := setupPedidoDB(t)
	c, _ := seedClienteAndProduto(t, db)
	repo := New(db)
	repo.Create(&entity.Pedido{ClienteID: c.ID, Status: "pendente", Total: 10.00})
	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestPedidoRepo_Delete(t *testing.T) {
	db := setupPedidoDB(t)
	c, _ := seedClienteAndProduto(t, db)
	repo := New(db)
	ped := &entity.Pedido{ClienteID: c.ID, Status: "pendente", Total: 10.00}
	repo.Create(ped)
	assert.NoError(t, repo.Delete(ped.ID))
	_, err := repo.FindByID(ped.ID)
	assert.Error(t, err)
}
