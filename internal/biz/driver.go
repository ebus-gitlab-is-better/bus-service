package biz

import "context"

type Driver struct {
	Id        *string
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Phone     *string `json:"number,omitempty"`
	BusNumber *string `json:"bus,omitempty"`
	Route     *string `json:"route,omitempty"`
}

type DriverRepo interface {
	GetDrivers(context.Context) ([]*Driver, error)
}

type DriverUseCase struct {
	repo DriverRepo
}

func NewDriverUseCase(repo DriverRepo) *DriverUseCase {
	return &DriverUseCase{repo: repo}
}

func (uc *DriverUseCase) GetDrivers(ctx context.Context) ([]*Driver, error) {
	return uc.repo.GetDrivers(ctx)
}
