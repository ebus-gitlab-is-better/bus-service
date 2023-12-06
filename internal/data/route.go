package data

import (
	"bus-service/internal/biz"
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	pq "github.com/lib/pq"
)

type Route struct {
	Id       uint32 `gorm:"primaryKey"`
	Number   string
	Path     string
	Time     pq.Float32Array `gorm:"type:double precision[]"`
	Lengths  pq.Float32Array `gorm:"type:double precision[]"`
	Stations []Stations      `gorm:"many2many:route_stations;"`
	Length   float32
}

func (m Route) modelToResponseWithoutStations() *biz.Route {
	return &biz.Route{
		Id:      m.Id,
		Number:  m.Number,
		Path:    m.Path,
		Time:    m.Time,
		Lengths: m.Lengths,
		Length:  m.Length,
	}
}

func (m Route) modelToResponse() *biz.Route {
	stations := make([]biz.Stations, 0)
	for _, station := range m.Stations {
		stations = append(stations, *station.modelToResponseWithoutRoute())
	}
	return &biz.Route{
		Id:       m.Id,
		Number:   m.Number,
		Path:     m.Path,
		Stations: stations,
		Time:     m.Time,
		Lengths:  m.Lengths,
		Length:   m.Length,
	}
}

type routeRepo struct {
	data   *Data
	logger *log.Helper
}

func NewRouterRepo(data *Data, logger log.Logger) biz.RouteRepo {
	return &routeRepo{data: data, logger: log.NewHelper(logger)}
}

// Create implements biz.RouteRepo.
func (r *routeRepo) Create(ctx context.Context, route *biz.Route) error {
	var routeDB Route
	routeDB.Number = route.Number
	routeDB.Path = route.Path
	routeDB.Time = route.Time
	routeDB.Lengths = route.Lengths
	stations := make([]Stations, 0)
	for _, station := range route.Stations {
		stations = append(stations, Stations{
			Name: station.Name,
			Lat:  station.Lat,
			Lon:  station.Lon,
		})
	}
	routeDB.Length = route.Length
	routeDB.Stations = stations
	if err := r.data.db.Create(&routeDB).Error; err != nil {
		return err
	}
	return nil
}

// Delete implements biz.RouteRepo.
func (r *routeRepo) Delete(ctx context.Context, id uint32) error {
	return r.data.db.Delete(&Route{}, id).Error
}

// GetById implements biz.RouteRepo.
func (r *routeRepo) GetById(ctx context.Context, id uint32) (*biz.Route, error) {
	var routeDB Route
	if err := r.data.db.Where(&Route{Id: id}).Find(&routeDB).Error; err != nil {
		return nil, err
	}
	return routeDB.modelToResponse(), nil
}

// List implements biz.RouteRepo.
func (r *routeRepo) List(context.Context) ([]*biz.Route, int64, error) {
	var routeDB []Route
	localDB := r.data.db.Model(&Route{})
	if err := localDB.Preload("Stations").Find(&routeDB).Error; err != nil {
		return nil, 0, err
	}
	var count int64
	localDB.Count(&count)
	route := make([]*biz.Route, 0)
	for _, b := range routeDB {
		route = append(route, b.modelToResponse())
	}
	fmt.Println(route[0].Stations)
	return route, count, nil
}

// Update implements biz.RouteRepo.
func (r *routeRepo) Update(ctx context.Context, route *biz.Route) error {
	var routeDB Route
	routeDB.Id = route.Id
	routeDB.Number = route.Number
	routeDB.Path = route.Path
	routeDB.Length = route.Length
	routeDB.Lengths = route.Lengths
	routeDB.Time = route.Time
	stations := make([]Stations, 0)
	for _, station := range route.Stations {
		stations = append(stations, Stations{
			ID:   station.ID,
			Name: station.Name,
			Lat:  station.Lat,
			Lon:  station.Lon,
		})
	}
	routeDB.Stations = stations
	if err := r.data.db.Save(&route).Error; err != nil {
		return err
	}
	return nil
}
