package model

import "time"

type Transaction struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"userId"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
