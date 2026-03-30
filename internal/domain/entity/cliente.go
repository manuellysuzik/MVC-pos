package entity

import "time"

type Cliente struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Nome      string    `json:"nome" gorm:"not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Telefone  string    `json:"telefone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
