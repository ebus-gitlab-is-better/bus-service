package biz

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Shift struct {
	Id        uint32
	StartTime time.Time
	EndDate   *time.Time
	DriverID  string
}

type ShiftRepo interface {
	Create(context.Context, *Shift) error
	Update(context.Context, *Shift) error
	GetByDriverID(context.Context, string) (*Shift, error)
}

type ShiftUseCase struct {
	repo ShiftRepo
}

func NewShiftUseCase(repo ShiftRepo) *ShiftUseCase {
	return &ShiftUseCase{repo: repo}
}

func (uc *ShiftUseCase) Create(ctx context.Context, shift *Shift) error {
	if data, err := uc.repo.GetByDriverID(ctx, shift.DriverID); data != nil || err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if data != nil {
			return errors.New("DRIVER_IN_DRIVE")
		}
	}
	return uc.repo.Create(ctx, shift)
}

func (uc *ShiftUseCase) Update(ctx context.Context, shift *Shift) error {
	return uc.repo.Update(ctx, shift)
}

func (uc *ShiftUseCase) GetByDriverID(ctx context.Context, driverId string) (*Shift, error) {
	return uc.repo.GetByDriverID(ctx, driverId)
}

func (uc *ShiftUseCase) GetHours(ctx context.Context, driverId string) (float64, error) {
	shift, err := uc.repo.GetByDriverID(ctx, driverId)
	if err != nil {
		return 0, err
	}
	endTime := time.Now()
	if shift.EndDate != nil {
		endTime = *shift.EndDate
	}

	duration := endTime.Sub(shift.StartTime)
	return duration.Hours(), nil
}
