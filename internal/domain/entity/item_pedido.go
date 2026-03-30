package entity

type ItemPedido struct {
	ID            uint     `json:"id" gorm:"primaryKey"`
	PedidoID      uint     `json:"pedido_id" gorm:"not null"`
	ProdutoID     uint     `json:"produto_id" gorm:"not null"`
	Produto       *Produto `json:"produto,omitempty" gorm:"foreignKey:ProdutoID"`
	Quantidade    int      `json:"quantidade" gorm:"not null"`
	PrecoUnitario float64  `json:"preco_unitario" gorm:"not null"`
}
