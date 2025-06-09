package dto

type PaginationResponse struct {
	RecordsFiltered int64       `json:"recordsFiltered"`
	Data            interface{} `json:"data"`
}
