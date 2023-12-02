package biz

import "context"

type Stations struct {
	ID     uint
	Name   string
	Lat    float64
	Lon    float64
	Routes []Route
}

type StationsPatch struct {
	ID     uint
	Name   *string
	Lat    *float64
	Lon    *float64
	Routes *[]*Route
}

type StationRepo interface {
	// Create(context.Context, *Stations) error
	Update(context.Context, *StationsPatch) error
	GetById(context.Context, uint32) (Stations, error)
	List(context.Context) ([]*Stations, int64, error)
	Delete(context.Context, uint32) error
}
