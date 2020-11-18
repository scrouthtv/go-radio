package tui

import "fmt"
import "time"

import "github.com/nsf/termbox-go"

import "github.com/scrouthtv/go-radio/stations"

type timelineScreen struct {
	station stations.Station
	day     time.Time
	cursor  Point
}

func (scr *timelineScreen) title() string {
	var name string
	var err error
	name, err = scr.station.GetName()
	if err != nil {
		name = fmt.Sprint(err)
	}
	return name
}

func (scr *timelineScreen) draw(lt Point, rb Point) {
	var row int

	var events []stations.Event
	var err error
	events, err = scr.station.DailyProgram(scr.day)
	row = centerprint(lt, Point{rb.X, lt.Y + 1}, coldef, coldef, scr.day.Format("< 02.01.2006 >"), false)
	if scr.cursor.Y == 0 {
		centerprint(lt, Point{rb.X, lt.Y + 1}, colcur, coldef, scr.day.Format("02.01.2006"), false)
	}
	if err != nil {
		tbprint(lt.X, lt.Y+row, coldef, coldef, "Error retrieving the program:")
		row++
		tbprint(lt.X, lt.Y+row, coldef, coldef, err)
		row++
	} else {
		tbprint(lt.X, lt.Y+row, coldef, coldef, "Retrieved ", len(events), " events")
		row++
	}
}

func (scr *timelineScreen) event(ev termbox.Event) bool {
	if scr.cursor.Y == 0 {
		if ev.Type == termbox.EventKey {
			if ev.Key == termbox.KeyArrowLeft {
				scr.day = scr.day.AddDate(0, 0, -1)
				return true
			} else if ev.Key == termbox.KeyArrowRight {
				scr.day = scr.day.AddDate(0, 0, 1)
				return true
			} else if ev.Key == termbox.KeyArrowDown {
				scr.cursor.Y++
			}
		}
	}
	return false
}

func (scr *timelineScreen) show() {
	currentScreen = scr
	scr.cursor = InvalidPoint
}

func (scr *timelineScreen) hide() {

}

func (scr *timelineScreen) focus() {
	focusScreen = scr
	scr.cursor = Point{0, 0}
}

func (scr *timelineScreen) unfocus() {
	scr.cursor = InvalidPoint
}
