package produtosvc

import "ecommerce-api/internal/domain/entity"

type ProdutoRepository interface {
	Create(produto *entity.Produto) error
	FindAll() ([]entity.Produto, error)
	FindByID(id uint) (*entity.Produto, error)
	FindByName(nome string) ([]entity.Produto, error)
	Count() (int64, error)
	Update(produto *entity.Produto) error
	Delete(id uint) error
}
