package unit

import (
	"math"
	"time"
)

const (
	Milisecond = "ms"
	Second     = "s"
	Minute     = "m"
	Hour       = "h"
	Day        = "D"
	Week       = "W"
	Year       = "Y"
)

type Sampling struct {
	Unit   string
	Period int
}

func GetUnitConversion(unit string) float64 {
	switch unit {
	case Milisecond:
		return float64(1.0)
	case Second:
		return GetUnitConversion(Milisecond) / 1000.0
	case Minute:
		return GetUnitConversion(Second) / 60.0
	case Hour:
		return GetUnitConversion(Minute) / 60.0
	case Day:
		return GetUnitConversion(Hour) / 24.0
	case Week:
		return GetUnitConversion(Day) / 7.0
	case Year:
		return GetUnitConversion(Day) / 365.0
	default:
		return 0.0
	}
}

func ConvertToUnit(date time.Time, unit string) float64 {
	switch unit {
	case Milisecond:
		return float64(date.UnixMilli())
	case Second:
		return ConvertToUnit(date, Milisecond) / 1000.0
	case Minute:
		return ConvertToUnit(date, Second) / 60.0
	case Hour:
		return ConvertToUnit(date, Minute) / 60.0
	case Day:
		return ConvertToUnit(date, Hour) / 24.0
	case Week:
		return ConvertToUnit(date, Day) / 7.0
	case Year:
		return ConvertToUnit(date, Day) / 365.0
	default:
		return 0.0
	}
}

func Sample(date time.Time, sampling Sampling) int {
	return int(math.Round(ConvertToUnit(date, sampling.Unit) / float64(sampling.Period)))
}
