package clientesvc

import "ecommerce-api/internal/domain/entity"

type ClienteRepository interface {
	Create(cliente *entity.Cliente) error
	FindAll() ([]entity.Cliente, error)
	FindByID(id uint) (*entity.Cliente, error)
	FindByName(nome string) ([]entity.Cliente, error)
	Count() (int64, error)
	Update(cliente *entity.Cliente) error
	Delete(id uint) error
}
