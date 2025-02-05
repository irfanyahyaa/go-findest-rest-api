package dto

type DashboardResponse struct {
	SuccessfulTransactionToday DashboardPagination[TransactionResponse] `json:"successfulTransactionToday"`
	AverageTransactionPerUser  []AverageTransactionAttr                 `json:"averageTransactionPerUser"`
	LatestTransaction          DashboardPagination[TransactionResponse] `json:"latestTransaction"`
}

type AverageTransactionAttr struct {
	UserId         uint    `json:"userId"`
	AvgTransaction float64 `json:"avgTransaction"`
}

type DashboardPagination[T any] struct {
	TotalRecords int `json:"totalRecords"`
	Transactions []T `json:"transactions"`
}
