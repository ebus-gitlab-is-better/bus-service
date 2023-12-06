package biz

import (
	"context"
	"encoding/json"
	"time"

	mapS "bus-service/api/map/v1"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/rabbitmq/amqp091-go"
)

type Route struct {
	Id       uint32
	Number   string
	Path     string
	Time     []float32
	Lengths  []float32
	Stations []Stations
	Length   float32
}

type Accident struct {
	Id        uint64     `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name"`
	Lat       float64    `json:"lat"`
	Lon       float64    `json:"lon"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

type RouteRepo interface {
	Create(context.Context, *Route) error
	Update(context.Context, *Route) error
	Delete(context.Context, uint32) error
	GetById(context.Context, uint32) (*Route, error)
	List(context.Context) ([]*Route, int64, error)
}

type RouteUseCase struct {
	repo      RouteRepo
	mapClient mapS.MapClient
	logger    *log.Helper
	rabbit    *RabbitData
}

func NewRouteUseCase(repo RouteRepo, logger log.Logger, mapClient mapS.MapClient, rabbit *RabbitData) *RouteUseCase {
	return &RouteUseCase{repo: repo, logger: log.NewHelper(logger), mapClient: mapClient, rabbit: rabbit}
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

type MessageDTO struct {
	Message string `json:"message"`
}

func (uc *RouteUseCase) NewAccident(ctx context.Context, accident *Accident) {
	routes, _, err := uc.List(context.TODO())
	if err != nil {
		return
	}
	for _, route := range routes {
		req, err := uc.mapClient.CheckPath(context.TODO(), &mapS.CheckPathRequest{
			Shape: route.Path,
			Point: &mapS.Point{
				Lat: float32(accident.Lat),
				Lon: float32(accident.Lon),
			},
		})
		if err != nil {
			return
		}
		jsonData, err := json.Marshal(&MessageDTO{
			Message: "Обнаружена дтп по маршруту " + route.Number + " извините за ожидание автобуса",
		})
		if err != nil {
			return
		}
		if req.IsValid {
			q, _ := uc.rabbit.Ch.QueueDeclare(
				"social", // name
				false,    // durable
				false,    // delete when unused
				false,    // exclusive
				false,    // no-wait
				nil,      // arguments
			)
			uc.rabbit.Ch.PublishWithContext(context.TODO(),
				"",
				q.Name,
				false,
				false,
				amqp091.Publishing{
					ContentType:  "application/json",
					Body:         jsonData,
					DeliveryMode: amqp091.Persistent,
				})
		}
	}
}
