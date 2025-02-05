package dto

import "time"

type TransactionCreate struct {
	UserID uint    `json:"userId"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}

type GetTransactionsQuery struct {
	UserID uint   `form:"userId"`
	Status string `form:"status"`
}

type Pagination[T any] struct {
	TotalRecords int `json:"totalRecords"`
	Data         []T `json:"data"`
}

type TransactionResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"userId"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TransactionUpdate struct {
	Status string `json:"status"`
}
