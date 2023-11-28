package profile

type Pagination struct {
	CurrentPage int `json:"currentPage"`
	PerPage     int `json:"perPage"`
	TotalPages  int `json:"totalPages"`
	TotalItems  int `json:"totalItems"`
}

func NewPagination(currentPage, perPage, totalItems int) Pagination {
	totalPages := totalItems / perPage
	if totalItems%perPage != 0 {
		totalPages += 1
	}

	return Pagination{
		CurrentPage: currentPage,
		PerPage:     perPage,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
	}
}
