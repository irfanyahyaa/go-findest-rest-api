package dto

type Pagination[T any] struct {
	TotalRecords int `json:"totalRecords"`
	Data         []T `json:"data"`
}
