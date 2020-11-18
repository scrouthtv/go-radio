package tui

import "github.com/nsf/termbox-go"

type sidebarScreen struct {
	cursor Point
}

func (scr *sidebarScreen) title() string {
	return "Menu"
}

func (scr *sidebarScreen) draw(lt Point, rb Point) {
	var i int
	var entry iScreen
	for i, entry = range availScreens {
		var title string = shrink(entry.title(), rb.X-lt.X)
		if scr.cursor != InvalidPoint && scr.cursor.Y == i {
			tbprint(lt.X, lt.Y+i, colcur, coldef, title)
		} else {
			tbprint(lt.X, lt.Y+i, coldef, coldef, title)
		}
	}
}

func (scr *sidebarScreen) event(ev termbox.Event) bool {
	if ev.Type == termbox.EventKey {
		if ev.Key == termbox.KeyArrowRight {
			scr.unfocus()
			currentScreen.focus()
			return true
		} else if ev.Key == termbox.KeyArrowDown {
			if scr.cursor.Y < len(availScreens)-1 {
				scr.cursor.Y = scr.cursor.Y + 1
				availScreens[scr.cursor.Y].show()
				return true
			}
		} else if ev.Key == termbox.KeyArrowUp {
			if scr.cursor.Y > 0 {
				scr.cursor.Y = scr.cursor.Y - 1
				availScreens[scr.cursor.Y].show()
				return true
			}
		}
	}
	return false
}

func (scr *sidebarScreen) show() {
	sidebar = scr
}

func (scr *sidebarScreen) hide() {

}

func (scr *sidebarScreen) focus() {
	focusScreen = scr
	scr.cursor = Point{0, 0}
}

func (scr *sidebarScreen) unfocus() {
	scr.cursor = InvalidPoint
}
