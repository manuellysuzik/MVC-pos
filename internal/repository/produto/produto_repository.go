package produtorepo

import (
	"errors"

	"ecommerce-api/internal/domain/entity"

	"gorm.io/gorm"
)

type ProdutoRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) ProdutoRepository {
	return ProdutoRepository{db: db}
}

func (r *ProdutoRepository) Create(p *entity.Produto) error {
	return r.db.Create(p).Error
}

func (r *ProdutoRepository) FindAll() ([]entity.Produto, error) {
	var produtos []entity.Produto
	return produtos, r.db.Find(&produtos).Error
}

func (r *ProdutoRepository) FindByID(id uint) (*entity.Produto, error) {
	var p entity.Produto
	if err := r.db.First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *ProdutoRepository) FindByName(nome string) ([]entity.Produto, error) {
	var produtos []entity.Produto
	return produtos, r.db.Where("nome LIKE ?", "%"+nome+"%").Find(&produtos).Error
}

func (r *ProdutoRepository) Count() (int64, error) {
	var count int64
	return count, r.db.Model(&entity.Produto{}).Count(&count).Error
}

func (r *ProdutoRepository) Update(p *entity.Produto) error {
	return r.db.Save(p).Error
}

func (r *ProdutoRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Produto{}, id).Error
}
