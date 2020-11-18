package tui

import "fmt"
import "time"

import "github.com/nsf/termbox-go"

import "github.com/scrouthtv/go-radio/stations"

type stationScreen struct {
	station  stations.Station
	day      time.Time
	cursor   Point
	events   []stations.Event // this could also be a pointer to the same array in timelineScreen
	err      error
	timeline *timelineScreen
}

func NewStationScreen(station stations.Station, day time.Time) stationScreen {
	return stationScreen{station, day, InvalidPoint, []stations.Event{}, nil,
		&timelineScreen{station, []stations.Event{}, InvalidPoint}}
}

func (scr *stationScreen) title() string {
	var name string
	var err error
	name, err = scr.station.GetName()
	if err != nil {
		name = fmt.Sprint(err)
	}
	return name
}

func (scr *stationScreen) draw(lt Point, rb Point) {
	var row int

	row = centerprint(lt, Point{rb.X, lt.Y + 1}, coldef, coldef, scr.day.Format("< 02.01.2006 >"), false)
	if scr.cursor.Y == 0 {
		centerprint(lt, Point{rb.X, lt.Y + 1}, colcur, coldef, scr.day.Format("02.01.2006"), false)
	}
	if scr.err != nil {
		tbprint(lt.X, lt.Y+row, coldef, coldef, "Error retrieving the program:")
		row++
		tbprint(lt.X, lt.Y+row, coldef, coldef, scr.err)
		row++
	} else {
		tbprint(lt.X, lt.Y+row, coldef, coldef, "Retrieved ", len(scr.events), " events")
		row++
	}

	scr.loadEvents()
	scr.timeline.draw(lt.Add(0, row+1), rb)
}

func (scr *stationScreen) loadEvents() {
	scr.events, scr.err = scr.station.DailyProgram(scr.day)
	scr.timeline.events = scr.events
}

func (scr *stationScreen) event(ev termbox.Event) bool {
	if scr.cursor.Y == 0 {
		if ev.Type == termbox.EventKey {
			if ev.Key == termbox.KeyArrowLeft {
				scr.day = scr.day.AddDate(0, 0, -1)
				scr.loadEvents()
				return true
			} else if ev.Key == termbox.KeyArrowRight {
				scr.day = scr.day.AddDate(0, 0, 1)
				scr.loadEvents()
				return true
			} else if ev.Key == termbox.KeyArrowDown {
				scr.cursor.Y++
				scr.timeline.focus()
				return true
			}
		}
	} else if scr.cursor.Y == 1 {
		if !scr.timeline.event(ev) {
			if ev.Type == termbox.EventKey {
				if ev.Key == termbox.KeyArrowUp {
					scr.timeline.unfocus()
					scr.cursor.Y--
				}
			}
		}
	}
	return false
}

func (scr *stationScreen) show() {
	currentScreen = scr
	scr.cursor = InvalidPoint
	scr.loadEvents()
}

func (scr *stationScreen) hide() {

}

func (scr *stationScreen) focus() {
	focusScreen = scr
	scr.timeline.unfocus()
	scr.cursor = Point{0, 0}
}

func (scr *stationScreen) unfocus() {
	scr.cursor = InvalidPoint
}
