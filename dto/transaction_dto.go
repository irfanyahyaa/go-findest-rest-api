package dto

import "time"

type TransactionDTO struct {
	UserID uint    `json:"userId"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}

type TransactionResponse struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"userId"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
}
