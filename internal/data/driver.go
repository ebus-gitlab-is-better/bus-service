package data

import (
	"bus-service/internal/biz"
	"context"
)

type Driver struct {
	Id        uint
	FirstName *string
	LastName  *string
	Phone     *string
}

type driverRepo struct {
	data *Data
}

func NewDriverRepo(data *Data) biz.DriverRepo {
	return &driverRepo{data: data}
}

const driverRole = "driver"

// GetDrivers implements biz.DriverRepo.
func (r *driverRepo) GetDrivers(ctx context.Context) ([]*biz.Driver, error) {
	kusers, err := r.data.keycloak.GetDrivers(driverRole)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0)
	for _, user := range kusers {
		ids = append(ids, *user.ID)
	}
	drivers := make([]*biz.Driver, 0)
	mapBus := map[string][]int{}
	buses, err := r.ListIn(ctx, ids)
	for i, r := range buses {
		if r.DriverID != nil {
			mapBus[*r.DriverID] = append(mapBus[*r.DriverID], i)
		}
	}
	for _, user := range kusers {
		dto := &biz.Driver{
			Id:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}
		if user.Attributes != nil {
			mapAttributes := *user.Attributes
			phone, ok := mapAttributes["phone"]
			if ok {
				dto.Phone = &phone[0]
			}
		}
		for _, index := range mapBus[*user.ID] {
			dto.BusNumber = &buses[index].Number
			if buses[index].Route != nil {
				dto.Route = &buses[index].Route.Number
			}
		}
		drivers = append(drivers, dto)
	}
	return drivers, err
}

func (r *driverRepo) ListIn(ctx context.Context, ids []string) ([]Bus, error) {
	var busDB []Bus
	localDB := r.data.db.Model(&Bus{})
	if err := localDB.Where("driver_id IN ?", ids).Find(&busDB).Error; err != nil {
		return nil, err
	}
	return busDB, nil
}

// {
// 	name: string
// 	rout: string
// 	bus: string
// 	status: string
// 	number: string
//   }
