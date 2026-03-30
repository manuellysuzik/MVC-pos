package clienterepo

import (
	"testing"

	"ecommerce-api/internal/domain/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupClienteDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&entity.Cliente{}))
	return db
}

func TestClienteRepo_Create(t *testing.T) {
	repo := New(setupClienteDB(t))
	c := &entity.Cliente{Nome: "João", Email: "joao@test.com"}
	assert.NoError(t, repo.Create(c))
	assert.NotZero(t, c.ID)
}

func TestClienteRepo_FindAll(t *testing.T) {
	repo := New(setupClienteDB(t))
	repo.Create(&entity.Cliente{Nome: "Ana", Email: "ana@test.com"})
	repo.Create(&entity.Cliente{Nome: "Bob", Email: "bob@test.com"})
	clientes, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, clientes, 2)
}

func TestClienteRepo_FindByID(t *testing.T) {
	repo := New(setupClienteDB(t))
	c := &entity.Cliente{Nome: "Carlos", Email: "carlos@test.com"}
	repo.Create(c)
	found, err := repo.FindByID(c.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Carlos", found.Nome)
}

func TestClienteRepo_FindByID_NotFound(t *testing.T) {
	repo := New(setupClienteDB(t))
	_, err := repo.FindByID(999)
	assert.Error(t, err)
}

func TestClienteRepo_FindByName(t *testing.T) {
	repo := New(setupClienteDB(t))
	repo.Create(&entity.Cliente{Nome: "Diego Costa", Email: "diego@test.com"})
	repo.Create(&entity.Cliente{Nome: "Diana Souza", Email: "diana@test.com"})
	repo.Create(&entity.Cliente{Nome: "Eduardo", Email: "edu@test.com"})
	results, err := repo.FindByName("Di")
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestClienteRepo_Count(t *testing.T) {
	repo := New(setupClienteDB(t))
	repo.Create(&entity.Cliente{Nome: "A", Email: "a@test.com"})
	repo.Create(&entity.Cliente{Nome: "B", Email: "b@test.com"})
	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestClienteRepo_Update(t *testing.T) {
	repo := New(setupClienteDB(t))
	c := &entity.Cliente{Nome: "Fernanda", Email: "fernanda@test.com"}
	repo.Create(c)
	c.Nome = "Fernanda Santos"
	assert.NoError(t, repo.Update(c))
	found, _ := repo.FindByID(c.ID)
	assert.Equal(t, "Fernanda Santos", found.Nome)
}

func TestClienteRepo_Delete(t *testing.T) {
	repo := New(setupClienteDB(t))
	c := &entity.Cliente{Nome: "Gabriel", Email: "gabriel@test.com"}
	repo.Create(c)
	assert.NoError(t, repo.Delete(c.ID))
	_, err := repo.FindByID(c.ID)
	assert.Error(t, err)
}
