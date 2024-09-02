package domainQuery

type FilmRepositoryQuery struct {
	SortField      string
	OrderBy        OrderDirection
	CurrentPage    int
	PageCount      int
	WithConnection []string
}

func NewFilmRepositoryQuery() *FilmRepositoryQuery {
	return &FilmRepositoryQuery{
		SortField:      "rate",
		OrderBy:        Asc,
		CurrentPage:    1,
		PageCount:      10,
		WithConnection: make([]string, 0, 4),
	}
}
