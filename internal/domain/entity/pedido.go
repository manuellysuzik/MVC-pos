package entity

import "time"

type Pedido struct {
	ID        uint         `json:"id" gorm:"primaryKey"`
	ClienteID uint         `json:"cliente_id" gorm:"not null"`
	Cliente   *Cliente     `json:"cliente,omitempty" gorm:"foreignKey:ClienteID"`
	Status    string       `json:"status" gorm:"default:'pendente'"`
	Total     float64      `json:"total"`
	Itens     []ItemPedido `json:"itens,omitempty" gorm:"foreignKey:PedidoID"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}
