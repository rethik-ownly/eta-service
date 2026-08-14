package repository

type Repository interface {
	GetETA(lat, lon float64) (float64, error)
}

type repositoryImpl struct {

}

func NewRepository() Repository {
	return &repositoryImpl{}
}

func (r *repositoryImpl) GetETA(lat, lon float64) (float64, error) {
	return 0, nil
}