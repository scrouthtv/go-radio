package tui

import "fmt"
import "time"

import "github.com/nsf/termbox-go"

import "github.com/scrouthtv/go-radio/stations"

type timelineScreen struct {
	station stations.Station
	day     time.Time
}

func (scr *timelineScreen) title() string {
	var name string
	var err error
	name, err = scr.station.GetName()
	if err != nil {
		name = fmt.Sprint(err)
	}
	return "Timeline for " + name
}

func (scr *timelineScreen) draw(lt Point, rb Point) {
	log("drawing for the ", scr.day.String())
	var events []stations.Event
	var err error
	events, err = scr.station.DailyProgram(scr.day)
	if err != nil {
		tbprint(lt.X, lt.Y, coldef, coldef, "Error retrieving the program:")
		tbprint(lt.X, lt.Y+1, coldef, coldef, err)
	} else {
		tbprint(lt.X, lt.Y, coldef, coldef, "Retrieved ", len(events), " events")
		tbprint(lt.X, lt.Y+1, coldef, coldef, "for the ", scr.day.Format("02.01.06"))
	}
}

func (scr *timelineScreen) event(ev termbox.Event) bool {
	if ev.Type == termbox.EventKey {
		if ev.Key == termbox.KeyArrowLeft {
			scr.day = scr.day.AddDate(0, 0, -1)
			return true
		} else if ev.Key == termbox.KeyArrowRight {
			scr.day = scr.day.AddDate(0, 0, 1)
			log("new date is ", scr.day.String())
			return true
		}
	}
	return false
}
