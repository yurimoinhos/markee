package common

// Page representa uma página de resultados paginados de forma genérica.
type Page[T any] struct {
	Items      []T `json:"items"`
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

// NewPage constrói uma Page[T] a partir dos itens, parâmetros de paginação e total de registros.
func NewPage[T any](items []T, page, pageSize, total int) Page[T] {
	totalPages := 0
	if pageSize > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}
	return Page[T]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}
