package clienterepo

import (
	"errors"

	"ecommerce-api/internal/domain/entity"

	"gorm.io/gorm"
)

type ClienteRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) ClienteRepository {
	return ClienteRepository{db: db}
}

func (r *ClienteRepository) Create(c *entity.Cliente) error {
	return r.db.Create(c).Error
}

func (r *ClienteRepository) FindAll() ([]entity.Cliente, error) {
	var clientes []entity.Cliente
	return clientes, r.db.Find(&clientes).Error
}

func (r *ClienteRepository) FindByID(id uint) (*entity.Cliente, error) {
	var c entity.Cliente
	if err := r.db.First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *ClienteRepository) FindByName(nome string) ([]entity.Cliente, error) {
	var clientes []entity.Cliente
	return clientes, r.db.Where("nome LIKE ?", "%"+nome+"%").Find(&clientes).Error
}

func (r *ClienteRepository) Count() (int64, error) {
	var count int64
	return count, r.db.Model(&entity.Cliente{}).Count(&count).Error
}

func (r *ClienteRepository) Update(c *entity.Cliente) error {
	return r.db.Save(c).Error
}

func (r *ClienteRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Cliente{}, id).Error
}
