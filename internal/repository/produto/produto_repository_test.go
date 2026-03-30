package produtorepo

import (
	"testing"

	"ecommerce-api/internal/domain/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupProdutoDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&entity.Produto{}))
	return db
}

func TestProdutoRepo_Create(t *testing.T) {
	repo := New(setupProdutoDB(t))
	p := &entity.Produto{Nome: "Notebook", Preco: 3500.00, Estoque: 10}
	assert.NoError(t, repo.Create(p))
	assert.NotZero(t, p.ID)
}

func TestProdutoRepo_FindAll(t *testing.T) {
	repo := New(setupProdutoDB(t))
	repo.Create(&entity.Produto{Nome: "Mouse", Preco: 80.00})
	repo.Create(&entity.Produto{Nome: "Teclado", Preco: 150.00})
	produtos, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, produtos, 2)
}

func TestProdutoRepo_FindByID(t *testing.T) {
	repo := New(setupProdutoDB(t))
	p := &entity.Produto{Nome: "Monitor", Preco: 1200.00}
	repo.Create(p)
	found, err := repo.FindByID(p.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Monitor", found.Nome)
}

func TestProdutoRepo_FindByID_NotFound(t *testing.T) {
	repo := New(setupProdutoDB(t))
	_, err := repo.FindByID(999)
	assert.Error(t, err)
}

func TestProdutoRepo_FindByName(t *testing.T) {
	repo := New(setupProdutoDB(t))
	repo.Create(&entity.Produto{Nome: "Cadeira Gamer", Preco: 800.00})
	repo.Create(&entity.Produto{Nome: "Cadeira Office", Preco: 500.00})
	repo.Create(&entity.Produto{Nome: "Mesa", Preco: 300.00})
	results, err := repo.FindByName("Cadeira")
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestProdutoRepo_Count(t *testing.T) {
	repo := New(setupProdutoDB(t))
	repo.Create(&entity.Produto{Nome: "A", Preco: 1.00})
	repo.Create(&entity.Produto{Nome: "B", Preco: 2.00})
	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestProdutoRepo_Update(t *testing.T) {
	repo := New(setupProdutoDB(t))
	p := &entity.Produto{Nome: "Headset", Preco: 200.00}
	repo.Create(p)
	p.Preco = 250.00
	assert.NoError(t, repo.Update(p))
	found, _ := repo.FindByID(p.ID)
	assert.Equal(t, 250.00, found.Preco)
}

func TestProdutoRepo_Delete(t *testing.T) {
	repo := New(setupProdutoDB(t))
	p := &entity.Produto{Nome: "Webcam", Preco: 300.00}
	repo.Create(p)
	assert.NoError(t, repo.Delete(p.ID))
	_, err := repo.FindByID(p.ID)
	assert.Error(t, err)
}
