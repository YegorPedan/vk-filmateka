package domainQuery

type ActorRepositoryQuery struct {
	CurrentPage    int
	PageCount      int
	WithConnection []string
}

func NewActorRepositoryQuery() *ActorRepositoryQuery {
	return &ActorRepositoryQuery{
		CurrentPage:    1,
		PageCount:      10,
		WithConnection: make([]string, 0, 4),
	}
}
