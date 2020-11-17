package stations

import "time"
import "fmt"

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
	End      time.Time
	Category string
}

// %n: name
// %i: info
// %s(15:08) start in the golang time format
// %e(15:08) end in the golang time format
// %c: category
func (ev *Event) Format(format string) string {
	return fmt.Sprintf("Start um %s - Name %s\n", ev.Start.Format("15:04"), ev.Name)
	//return fmt.Sprintf("Start um %s - Name %s\n%s\n", ev.Start.Format("15:04"), ev.Name, ev.Info)
}
