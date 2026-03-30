package entity

import "time"

type Produto struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Nome      string    `json:"nome" gorm:"not null"`
	Descricao string    `json:"descricao"`
	Preco     float64   `json:"preco" gorm:"not null"`
	Estoque   int       `json:"estoque" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
