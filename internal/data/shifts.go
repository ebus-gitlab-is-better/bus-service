package data

import (
	"bus-service/internal/biz"
	"context"
	"time"
)

type Shift struct {
	Id        uint32 `gorm:"primaryKey"`
	StartTime time.Time
	EndDate   *time.Time
	DriverID  string
}

func (m Shift) modelToResponse() *biz.Shift {
	return &biz.Shift{
		Id:        m.Id,
		StartTime: m.StartTime,
		EndDate:   m.EndDate,
		DriverID:  m.DriverID,
	}
}

type shiftRepo struct {
	data *Data
}

func NewShiftRepo(data *Data) biz.ShiftRepo {
	return &shiftRepo{data: data}
}

// Create implements biz.ShiftRepo.
func (r *shiftRepo) Create(ctx context.Context, shift *biz.Shift) error {
	shiftDB := Shift{}
	shiftDB.StartTime = shift.StartTime
	shiftDB.EndDate = shift.EndDate
	shiftDB.DriverID = shift.DriverID
	if err := r.data.db.Create(&shiftDB).Error; err != nil {
		return err
	}
	return nil
}

// GetById implements biz.ShiftRepo.
func (r *shiftRepo) GetByDriverID(ctx context.Context, driverId string) (*biz.Shift, error) {
	var shiftDB Shift
	if err := r.data.db.Where(&Shift{DriverID: driverId, EndDate: nil}).Order("start_time DESC").First(&shiftDB).Error; err != nil {
		return nil, err
	}
	return shiftDB.modelToResponse(), nil
}

// Update implements biz.ShiftRepo.
func (r *shiftRepo) Update(ctx context.Context, shift *biz.Shift) error {
	shiftDB := Shift{}
	shiftDB.StartTime = shift.StartTime
	shiftDB.EndDate = shift.EndDate
	shiftDB.Id = shift.Id
	shiftDB.DriverID = shift.DriverID
	if err := r.data.db.Save(&shiftDB).Error; err != nil {
		return err
	}
	return nil
}
