package tui

import "fmt"
import "time"
import "math"

import "github.com/nsf/termbox-go"

import "github.com/scrouthtv/go-radio/stations"

type timelineScreen struct {
	station stations.Station
	day     time.Time
	cursor  Point
	events  []stations.Event
	err     error
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

	// draw timeline:
	tbprint(lt.X, lt.Y+row+1, coldef, coldef, "<")
	tbprint(rb.X, lt.Y+row+1, coldef, coldef, ">")
	const minsPerSpace int = 15
	var x, w, i int
	var startMinute, lastMinute = 60 * 9, 60 * 9 // start at 9:00
	var events []stations.Event = SortEventsByStart(scr.events)
	var event stations.Event
	var innerRow int = 0
	var eMin int // event minute
	for i, event = range events {
		eMin = timeAsMinutes(event.Start)
		if eMin > lastMinute {
			x = (eMin - startMinute) / minsPerSpace
			w = int(math.Ceil(event.End.Sub(event.Start).Minutes())) / minsPerSpace
			log("drawing ", i, ": ", event.Name, " to ", lt.X+2+x, "x", lt.Y+row+innerRow)
			if scr.cursor.Y == 1 && scr.cursor.X == i {
				log("highlighted")
				tbprint(lt.X+2+x, lt.Y+row+innerRow, colcur, coldef, event.Name[0:2])
			} else {
				tbprint(lt.X+2+x, lt.Y+row+innerRow, coldef, coldef, event.Name[0:2])
			}
			innerRow = (innerRow + 1) % 3
			lastMinute = eMin + w
		}
	}
}

func (scr *timelineScreen) event(ev termbox.Event) bool {
	if scr.cursor.Y == 0 {
		if ev.Type == termbox.EventKey {
			if ev.Key == termbox.KeyArrowLeft {
				scr.day = scr.day.AddDate(0, 0, -1)
				scr.events, scr.err = scr.station.DailyProgram(scr.day)
				return true
			} else if ev.Key == termbox.KeyArrowRight {
				scr.day = scr.day.AddDate(0, 0, 1)
				scr.events, scr.err = scr.station.DailyProgram(scr.day)
				return true
			} else if ev.Key == termbox.KeyArrowDown {
				scr.cursor.Y++
				return true
			}
		}
	} else if scr.cursor.Y == 1 {
		if ev.Type == termbox.EventKey {
			if ev.Key == termbox.KeyArrowLeft {
				if scr.cursor.X > 0 {
					scr.cursor.X--
					return true
				}
			} else if ev.Key == termbox.KeyArrowRight {
				scr.cursor.X++
				return true
			} else if ev.Key == termbox.KeyArrowUp {
				scr.cursor.Y--
				return true
			}
		}
	}
	return false
}

func (scr *timelineScreen) show() {
	currentScreen = scr
	scr.cursor = InvalidPoint
	scr.events, scr.err = scr.station.DailyProgram(scr.day)
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
