package repositories

func NewRepository() Repository {
	repo := make(MapRepository, 0)
	return &repo
}
