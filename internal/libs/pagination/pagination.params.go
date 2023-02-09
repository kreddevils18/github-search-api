package pagination

import (
	"github.com/kien-hoangtrung/github-repository/internal/libs/pagination/dto"
	"math"
)

const (
	defaultLimit = 10
	defaultPage  = 1
	defaultSkip  = 0
)

type Params interface {
	dto.PaginationRequestInterface
}

func PaginationParamsTransform[T Params](req T) T {
	limit := req.GetLimit()
	page := req.GetPage()
	skip := req.GetSkip()

	if limit == 0 {
		req.SetLimit(defaultLimit)
	}

	if skip == 0 {
		if page != 0 {
			req.SetSkip((page - 1) * limit)
		} else {
			req.SetPage(defaultPage)
			req.SetSkip(defaultSkip)
		}
	} else {
		req.SetPage(int64(math.Floor(float64(skip) / float64(limit))))
	}

	return req
}
