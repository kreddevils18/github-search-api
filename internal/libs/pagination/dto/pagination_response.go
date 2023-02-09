package dto

type PaginationResponse struct {
	CurrentPage    int64       `json:"current_page"`
	SkippedRecords int64       `json:"skipped_records"`
	TotalRecords   int64       `json:"total_records"`
	TotalPages     int64       `json:"total_pages"`
	HasNext        bool        `json:"has_next"`
	Data           interface{} `json:"data"`
}
