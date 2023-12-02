package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type Route struct {
	Id       uint32
	Number   string
	Path     string
	Stations []Stations
}

type RouteRepo interface {
	Create(context.Context, *Route) error
	Update(context.Context, *Route) error
	Delete(context.Context, uint32) error
	GetById(context.Context, uint32) (*Route, error)
	List(context.Context) ([]*Route, int64, error)
}

type RouteUseCase struct {
	repo   RouteRepo
	logger *log.Helper
}

func NewRouteUseCase(repo RouteRepo, logger log.Logger) *RouteUseCase {
	return &RouteUseCase{repo: repo, logger: log.NewHelper(logger)}
}

func (uc *RouteUseCase) Create(ctx context.Context, route *Route) error {
	return uc.repo.Create(ctx, route)
}

func (uc *RouteUseCase) Update(ctx context.Context, route *Route) error {
	return uc.repo.Update(ctx, route)
}

func (uc *RouteUseCase) Delete(ctx context.Context, id uint32) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *RouteUseCase) GetById(ctx context.Context, id uint32) (*Route, error) {
	return uc.repo.GetById(ctx, id)
}

func (uc *RouteUseCase) List(ctx context.Context) ([]*Route, int64, error) {
	return uc.repo.List(ctx)
}
