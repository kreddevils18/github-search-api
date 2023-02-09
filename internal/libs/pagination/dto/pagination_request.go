package dto

type OrderType int

type PaginationRequestInterface interface {
	SetSkip(int64)
	GetSkip() int64
	SetPage(int64)
	GetPage() int64
	SetLimit(int64)
	GetLimit() int64
}

type PaginationRequest struct {
	Skip  int64 `json:"skip" query:"skip" validate:"min=0"`
	Page  int64 `json:"page" query:"page" validate:"min=0"`
	Limit int64 `json:"limit" query:"limit" validate:"min=0,max=30"`
}

func (p *PaginationRequest) SetSkip(skip int64) {
	p.Skip = skip
}

func (p PaginationRequest) GetSkip() int64 {
	return p.Skip
}

func (p *PaginationRequest) SetPage(page int64) {
	p.Page = page
}

func (p PaginationRequest) GetPage() int64 {
	return p.Page
}

func (p *PaginationRequest) SetLimit(limit int64) {
	p.Limit = limit
}

func (p PaginationRequest) GetLimit() int64 {
	return p.Limit
}
