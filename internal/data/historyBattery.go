package data

import "time"

type HistoryBattery struct {
	Id            uint
	Date          time.Time
	RouteDistance float32
	BatterUsage   float32
	RouteID       uint32
}
