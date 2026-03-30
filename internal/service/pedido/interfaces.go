package pedidosvc

import "ecommerce-api/internal/domain/entity"

type PedidoRepository interface {
	Create(pedido *entity.Pedido) error
	FindAll() ([]entity.Pedido, error)
	FindByID(id uint) (*entity.Pedido, error)
	FindByStatus(status string) ([]entity.Pedido, error)
	FindByClienteID(clienteID uint) ([]entity.Pedido, error)
	Count() (int64, error)
	Update(pedido *entity.Pedido) error
	Delete(id uint) error
}
