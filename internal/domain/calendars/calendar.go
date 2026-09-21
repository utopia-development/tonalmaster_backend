package calendars

import (
	"errors"
	"time"
)

var ErrUnsupportedSystem = errors.New("unsupported calendar system")

type Result struct {
	System   string
	Date     time.Time
	JDN      int64
	Trecena  int
	Sign     string
	DayNumber int
	NightLord string
}

type System interface {
	ID() string
	Convert(date time.Time) (Result, error)
}

func Convert(system System, date time.Time) (Result, error) {
	return system.Convert(date)
}
