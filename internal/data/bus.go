package data

import (
	"bus-service/internal/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type Bus struct {
	Id       uint32 `gorm:"primaryKey"`
	RouteID  uint32
	DriverID string
}

type busRepo struct {
	data   *Data
	logger *log.Helper
}

func NewBusRepo(data *Data, logger log.Logger) biz.BusRepo {
	return &busRepo{data: data, logger: log.NewHelper(logger)}
}

// Create implements biz.BusRepo.
func (r *busRepo) Create(ctx context.Context, bus *biz.BusDTO) error {
	var busDB Bus
	busDB.RouteID = bus.RouteID
	busDB.DriverID = bus.DriverID
	if err := r.data.db.Create(&busDB).Error; err != nil {
		return err
	}
	return nil
}

// Delete implements biz.BusRepo.
func (r *busRepo) Delete(ctx context.Context, id uint32) error {
	return r.data.db.Delete(&Bus{}, id).Error
}

// GetById implements biz.BusRepo.
func (r *busRepo) GetById(ctx context.Context, id uint32) (*biz.Bus, error) {
	var busDB Bus
	if err := r.data.db.Where(&Bus{Id: id}).Find(&busDB).Error; err != nil {
		return nil, err
	}
	return r.modelToResponse(busDB), nil
}

// List implements biz.BusRepo.
func (r *busRepo) List(ctx context.Context) ([]*biz.Bus, int64, error) {
	var busDB []Bus
	localDB := r.data.db.Model(&Bus{})
	if err := localDB.Find(&busDB).Error; err != nil {
		return nil, 0, err
	}
	var count int64
	localDB.Count(&count)
	bus := make([]*biz.Bus, 0)
	for _, b := range busDB {
		bus = append(bus, r.modelToResponse(b))
	}
	return bus, count, nil
}

// Update implements biz.BusRepo.
func (r *busRepo) Update(ctx context.Context, bus *biz.BusDTO) error {
	var busDB Bus
	busDB.RouteID = bus.RouteID
	busDB.DriverID = bus.DriverID
	busDB.Id = bus.Id
	if err := r.data.db.Save(&busDB).Error; err != nil {
		return err
	}
	return nil
}

func (r *busRepo) modelToResponse(b Bus) *biz.Bus {
	user, _ := r.data.keycloak.GetUserByID(b.DriverID)
	//ignore err
	return &biz.Bus{
		Id:      b.Id,
		RouteID: b.RouteID,
		Driver: biz.BusUser{
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Id:        user.ID,
		},
	}
}
