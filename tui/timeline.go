package tui

import "fmt"

import "github.com/nsf/termbox-go"

import "github.com/scrouthtv/go-radio/stations"

type timelineScreen struct {
	station   stations.Station
	events    []stations.Event
	cursor    Point
	selected  *stations.Event
	highlight func(station *stations.Station, ev *stations.Event)
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
	tbprint(lt.X, lt.Y, coldef, coldef, "<")
	tbprint(rb.X, lt.Y, coldef, coldef, ">")
	const minsPerSpace int = 5
	const headerInterval int = 60 // in minutes
	const headerWidth int = 5
	const startMinute = 60 * 9 // start at 9:00

	var x, w, count int
	var lastMinute int = startMinute

	// draw the header:
	var time, h, m int
	for time = startMinute; x < rb.X; time += headerInterval {
		h, m = minutesToHM(time)
		tbprint(lt.X+2+x, lt.Y, coldef, coldef, fmt.Sprintf("%02d:%02d", h, m))
		x += headerInterval / minsPerSpace
	}

	// draw the actual timeline:
	var events []stations.Event = SortEventsByStart(scr.events)
	var event, sEv stations.Event
	var row int = 1
	var eMin int // event minute
	x = 0
	for _, event = range events { // for each event:
		eMin = timeAsMinutes(event.Start)
		if eMin > lastMinute { // if it is later than the start minute
			x = (eMin - startMinute) / minsPerSpace // calculate the start
			//w = int(math.Ceil(event.End.Sub(event.Start).Minutes())) / minsPerSpace
			w = 6
			if len(event.Name) < w {
				w = len(event.Name)
			}
			if lt.X+2+x+w > rb.X+2 {
				w = rb.X - lt.X - 2 - x - 2
			}
			if w <= 0 {
				continue
			}
			if scr.cursor.Y == 0 && scr.cursor.X == count {
				tbprint(lt.X+2+x, lt.Y+row, colcur, coldef, event.Name[0:w])
				sEv = event
				scr.selected = &sEv // this is a very hacky implementation that should be fixed.
				// the basic problem is that 'event' itself will get overridden in the next loop run
				// so I copy it to sEv and select the pointer to that one. Instead, I should loop
				// over the pointers to the elements in the events array
				scr.highlight(&scr.station, &sEv)
			} else {
				tbprint(lt.X+2+x, lt.Y+row, coldef, coldef, event.Name[0:w])
			}
			//row = (row + 1) % 3
			lastMinute = eMin + w
			count++
		}
	}

}

func (scr *timelineScreen) event(ev termbox.Event) bool {
	if ev.Type == termbox.EventKey {
		if ev.Key == termbox.KeyArrowRight {
			scr.cursor.X++
			return true
		} else if ev.Key == termbox.KeyArrowLeft {
			scr.cursor.X--
			return true
		} // up & down is handled by the station screen
	}
	return false
}

func (scr *timelineScreen) show() {
}

func (scr *timelineScreen) hide() {
}

func (scr *timelineScreen) focus() {
	scr.cursor = Point{0, 0}
}

func (scr *timelineScreen) unfocus() {
	scr.cursor = InvalidPoint
}
