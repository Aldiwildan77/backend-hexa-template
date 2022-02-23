package pagination

import "math"

type Pagination struct {
	Limit  int
	Page   int
	Total  int
	From   int
	To     int
	Last   int
	Offset int
}

func (p *Pagination) ValidatePagination() {
	if p.Limit == 0 {
		p.Limit = 10
	}

	if p.Page == 0 {
		p.Page = 1
	}

	p.From = CalculateFrom(p.Limit, p.Page)
	p.To = CalculateTo(p.Limit, p.Page)
	p.Offset = CalculateOffset(p.Limit, p.Page)
	p.Last = CalculateLastPage(p.Total, p.Limit)
}

func CalculateOffset(limit, page int) int {
	return (page - 1) * limit
}

func CalculateFrom(limit, page int) int {
	offset := CalculateOffset(limit, page)

	if offset > 1 {
		return offset + 1
	}

	return 1
}

func CalculateTo(limit, page int) int {
	from := CalculateFrom(limit, page)

	return from + limit - 1
}

func CalculateLastPage(total, limit int) int {
	result := math.Ceil(float64(total) / float64(limit))

	if int(result) == 0 {
		return 1
	}

	return int(result)
}
