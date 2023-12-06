package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type Bus struct {
	Id      uint32
	RouteID *uint32
	Route   *Route
	Driver  BusUser
	Number  string
	Status  string
}

type BusDTO struct {
	Id       uint32
	RouteID  *uint32
	DriverID *string
	Number   string
	Status   string
}

type BusUser struct {
	Id        *string
	Username  *string
	FirstName *string
	LastName  *string
	Email     *string
}

type BusRepo interface {
	Create(context.Context, *BusDTO) error
	Update(context.Context, *BusDTO) error
	GetById(context.Context, uint32) (*Bus, error)
	List(context.Context) ([]*Bus, int64, error)
	Delete(context.Context, uint32) error
	GetActiveBus(context.Context) ([]*Bus, error)
}

type BusUseCase struct {
	repo   BusRepo
	logger *log.Helper
}

func NewBusUseCase(repo BusRepo, logger log.Logger) *BusUseCase {
	return &BusUseCase{repo: repo, logger: log.NewHelper(logger)}
}

func (uc *BusUseCase) Create(ctx context.Context, bus *BusDTO) error {
	return uc.repo.Create(ctx, bus)
}

func (uc *BusUseCase) Update(ctx context.Context, bus *BusDTO) error {
	return uc.repo.Update(ctx, bus)
}

func (uc *BusUseCase) GetById(ctx context.Context, id uint32) (*Bus, error) {
	return uc.repo.GetById(context.TODO(), id)
}

func (uc *BusUseCase) Delete(ctx context.Context, id uint32) error {
	return uc.repo.Delete(context.TODO(), id)
}

func (uc *BusUseCase) List(ctx context.Context) ([]*Bus, int64, error) {
	return uc.repo.List(ctx)
}
