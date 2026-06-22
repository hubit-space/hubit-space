package model

type PaginationResult struct {
	Status       int   `json:"status"`
	CurrentPage  int   `json:"current_page"`
	Data         any   `json:"data"`
	From         int   `json:"from"`
	To           int   `json:"to"`
	PerPage      int   `json:"per_page"`
	Total        int64 `json:"total"`
	TotalPages   int   `json:"total_pages"`
	PreviousPage int   `json:"prev,omitempty"`
	NextPage     int   `json:"next,omitempty"`
}
