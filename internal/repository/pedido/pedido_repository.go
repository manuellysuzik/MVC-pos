package pedidorepo

import (
	"errors"

	"ecommerce-api/internal/domain/entity"

	"gorm.io/gorm"
)

type PedidoRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) PedidoRepository {
	return PedidoRepository{db: db}
}

func (r *PedidoRepository) Create(p *entity.Pedido) error {
	return r.db.Create(p).Error
}

func (r *PedidoRepository) FindAll() ([]entity.Pedido, error) {
	var pedidos []entity.Pedido
	return pedidos, r.db.Preload("Itens").Preload("Cliente").Find(&pedidos).Error
}

func (r *PedidoRepository) FindByID(id uint) (*entity.Pedido, error) {
	var p entity.Pedido
	if err := r.db.Preload("Itens").Preload("Cliente").First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *PedidoRepository) FindByStatus(status string) ([]entity.Pedido, error) {
	var pedidos []entity.Pedido
	return pedidos, r.db.Preload("Itens").Where("status = ?", status).Find(&pedidos).Error
}

func (r *PedidoRepository) FindByClienteID(clienteID uint) ([]entity.Pedido, error) {
	var pedidos []entity.Pedido
	return pedidos, r.db.Preload("Itens").Where("cliente_id = ?", clienteID).Find(&pedidos).Error
}

func (r *PedidoRepository) Count() (int64, error) {
	var count int64
	return count, r.db.Model(&entity.Pedido{}).Count(&count).Error
}

func (r *PedidoRepository) Update(p *entity.Pedido) error {
	return r.db.Save(p).Error
}

func (r *PedidoRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Pedido{}, id).Error
}
