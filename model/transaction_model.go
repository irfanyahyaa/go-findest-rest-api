package model

import "time"

type Transaction struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"userId"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	IsDeleted bool      `json:"isDeleted"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
