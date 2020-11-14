package stations

import "time"

type Station interface {
	GetName() (string, error)
	GetURL() (string, error)
	Program() ([]Event, error)
	DailyProgram(day time.Time) ([]Event, error)
}

type Event struct {
	Name     string
	Info     string
	Start    time.Time
	Duration time.Duration
	Category string
}
