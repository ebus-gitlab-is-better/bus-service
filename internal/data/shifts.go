package data

import "time"

type Shift struct {
	Id        uint `gorm:"primaryKey"`
	StartTime time.Time
	EndDate   *time.Time
	DriverID  uint
}
