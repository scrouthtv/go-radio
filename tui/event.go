package tui

import "net/url"

import "github.com/nsf/termbox-go"

import "github.com/scrouthtv/go-radio/stations"

type eventScreen struct {
	ev *stations.Event
}

func (scr *eventScreen) title() string {
	if scr.ev == nil {
		return "Event"
	} else {
		return scr.ev.Name
	}
}

func (scr *eventScreen) draw(lt Point, rb Point) {
	if scr.ev == nil {
		return
	}

	//var ltX = lt.X
	var rbX = rb.X

	var row int = 0

	// name & time full centered
	row += centerprint(lt.Add(0, row), rb, colimp, colred, shrink(scr.ev.Name, rbX-2), true)
	var t string = scr.ev.Start.Format("15:04") + " - " + scr.ev.End.Format("15:04")
	row += centerprint(lt.Add(0, row), rb, coldef, coldef, t, false)

	// info text on the left side
	rb.X = (rbX-lt.X)/2 + lt.X
	boxprint(lt.Add(0, row), rb, colgry, coldef, scr.ev.Info)

	// more information on the right side
	lt.X = rb.X
	rb.X = rbX
	tbprint(lt.X, lt.Y+row, coldef, coldef, "Appears on ", scr.ev.Category)
	row++
	var emptyUrl *url.URL
	emptyUrl, _ = url.Parse("")
	if scr.ev.CatUrl != *emptyUrl {
		tbprint(lt.X, lt.Y+row, coldef, coldef, " (", scr.ev.CatUrl.String(), ")")
		row++
	}
	//if scr.ev.Url != nil {
	tbprint(lt.X, lt.Y+row, coldef, coldef, "More info via ", scr.ev.Url.String())
	row++
	//}
}

func (scr *eventScreen) event(ev termbox.Event) bool {
	return false
}

func (scr *eventScreen) show() {

}

func (scr *eventScreen) hide() {

}

func (scr *eventScreen) focus() {

}

func (scr *eventScreen) unfocus() {

}
