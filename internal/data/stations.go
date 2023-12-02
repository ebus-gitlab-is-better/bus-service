package data

import (
	"bus-service/internal/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type Stations struct {
	ID     uint `gorm:"primaryKey"`
	Name   string
	Lat    float64
	Lon    float64
	Routes []Route `gorm:"many2many:route_stations;"`
}

func (m Stations) modelToResponse() *biz.Stations {
	routes := make([]biz.Route, 0)
	for _, route := range m.Routes {
		routes = append(routes, *route.modelToResponse())
	}
	return &biz.Stations{
		ID:     m.ID,
		Name:   m.Name,
		Lat:    m.Lat,
		Lon:    m.Lon,
		Routes: routes,
	}
}

func (m Stations) modelToResponseWithoutRoute() *biz.Stations {
	return &biz.Stations{
		ID:   m.ID,
		Name: m.Name,
		Lat:  m.Lat,
		Lon:  m.Lon,
	}
}

type stationsRepo struct {
	data   *Data
	logger *log.Helper
}

func NewStationsRepo(data *Data, logger log.Logger) biz.StationRepo {
	return &stationsRepo{data: data, logger: log.NewHelper(logger)}
}

// Create implements biz.StationRepo.
// func (r *stationsRepo) Create(ctx context.Context, station *biz.Stations) error {
// 	stationDB :=
// }

// Delete implements biz.StationRepo.
func (*stationsRepo) Delete(context.Context, uint32) error {
	panic("unimplemented")
}

// GetById implements biz.StationRepo.
func (*stationsRepo) GetById(context.Context, uint32) (biz.Stations, error) {
	panic("unimplemented")
}

// List implements biz.StationRepo.
func (*stationsRepo) List(context.Context) ([]*biz.Stations, int64, error) {
	panic("unimplemented")
}

// Update implements biz.StationRepo.
func (*stationsRepo) Update(context.Context, *biz.StationsPatch) error {
	panic("unimplemented")
}
