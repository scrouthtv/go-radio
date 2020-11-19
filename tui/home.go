package tui

import "time"

import "github.com/nsf/termbox-go"

type homeScreen struct {
}

func (scr *homeScreen) title() string {
	return "Home"
}

func (scr *homeScreen) draw(lt Point, rb Point) {
	var hour int = time.Now().Hour()
	if hour <= 9 {
		tbprint(lt.X, lt.Y, coldef, coldef, "Good Morning")
	} else if hour <= 14 {
		tbprint(lt.X, lt.Y, coldef, coldef, "G'day")
	} else if hour <= 18 {
		tbprint(lt.X, lt.Y, coldef, coldef, "Good Afternoon")
	} else {
		tbprint(lt.X, lt.Y, coldef, coldef, "Good Evening")
	}
	/*var ls string = "Here is a long string that will require line breaks in between because I put some many words there that it is simply unbelievable how many words can and can't fit in one line I did this just to test the split functionality. Also there should be unbelievably schtupid long words like thisonewherethereisnospacesbecauseitisasupercomplicatedwordthatissuporlongandrequiressplittingalongmultiplelineswhichIdidimplementbutIwanttotestit also there are words after this super large one that should be printed on their own lines as well."
	var offset int = boxprint(lt.Add(0, 1), rb, coldef, coldef, ls) + 1
	ls = "This is centered text."
	offset += centerprint(lt.Add(0, offset), Point{lt.X + 10, rb.Y}, colcur, coldef, ls, true)
	centerprint(lt.Add(0, offset), Point{lt.X + 10, rb.Y}, colcur, coldef, ls, false)*/
}

func (scr *homeScreen) event(ev termbox.Event) bool {
	if ev.Type == termbox.EventKey && ev.Key == termbox.KeyArrowLeft {
		sidebar.focus()
		return true
	}
	return false
}

func (scr *homeScreen) show() {
	currentScreen = scr
}

func (scr *homeScreen) hide() {
}

func (scr *homeScreen) focus() {
	focusScreen = scr
}

func (scr *homeScreen) unfocus() {
}
