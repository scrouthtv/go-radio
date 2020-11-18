package tui

import "fmt"
import "time"
import "strings"
import "math"

import "github.com/nsf/termbox-go"

import "github.com/scrouthtv/go-radio/stations"

type iScreen interface {
	title() string
	draw(lt Point, rb Point)
	event(ev termbox.Event) bool
	show()
	hide()
	focus()
	unfocus()
}

type Point struct {
	X int
	Y int
}

var InvalidPoint Point = Point{-1, -1}

type Area struct {
	Lt Point
	Rb Point
}

var Running bool = false

const coldef termbox.Attribute = termbox.ColorDefault
const colcur termbox.Attribute = coldef | termbox.AttrReverse
const colsel termbox.Attribute = coldef | termbox.AttrBold

var currentScreen iScreen
var focusScreen iScreen
var sidebar *sidebarScreen

var availScreens []iScreen = []iScreen{}

// this might not work since iScreens don't implement any comparation functionality
var overlays map[iScreen]Area = make(map[iScreen]Area)

var sepWidth int

const padding int = 1

func TuiLoop() error {
	Running = true
	var err error
	err = termbox.Init()
	if err != nil {
		return err
	}

	var scr stationScreen = NewStationScreen(stations.Deutschlandfunk, time.Now())
	availScreens = append(availScreens, &homeScreen{}, &scr)

	availScreens[1].show()
	currentScreen.focus()
	var sidebar iScreen = &sidebarScreen{InvalidPoint}
	sidebar.show()

	var w, h int = termbox.Size()
	sepWidth = w / 7
	overlays[sidebar] = Area{Point{0, 0}, Point{sepWidth - 1, h}}

	// TODO mouse support?, resizing

	var ev termbox.Event
	for {
		redraw()
		ev = termbox.PollEvent()
		if !focusScreen.event(ev) {
			tbprint(10, 10, coldef, coldef, "event was not consumed")
			if ev.Type == termbox.EventKey && ev.Ch == 'q' {
				break
			}
		}
	}
	termbox.Close()
	Running = false
	return nil
}

func tbprint(x int, y int, fg termbox.Attribute, bg termbox.Attribute, a ...interface{}) {
	var msg string = fmt.Sprint(a...)
	var i int
	var char rune
	for i, char = range []rune(msg) {
		termbox.SetCell(x+i, y, char, fg, bg)
	}
}

func boxprint(lt Point, rb Point, fg termbox.Attribute, bg termbox.Attribute, msg string) int {
	var width int = rb.X - lt.X
	var lines []string = softwrap(msg, width)
	var rows int = len(lines)
	if rows > rb.Y-lt.Y {
		rows = rb.Y - lt.Y
	}
	var row int
	for row = 0; row < rows; row++ {
		tbprint(lt.X, lt.Y+row, fg, bg, lines[row])
	}
	return rows
}

func centerprint(lt Point, rb Point, fg termbox.Attribute, bg termbox.Attribute, msg string, surroundWithAttr bool) int {
	var width int = rb.X - lt.X
	var lines []string = softwrap(msg, width)
	var rows int = len(lines)
	if rows > rb.Y-lt.Y {
		rows = rb.Y - lt.Y
	}
	var row int
	var pfx, sfx int
	var line []rune

	for row = 0; row < rows; row++ {
		line = []rune(lines[row])
		if IsFullRune(line[len(line)-1]) {
			pfx = int(math.Floor(float64(width-len(lines[row])) / 2.0))
			sfx = int(math.Ceil(float64(width-len(lines[row])) / 2.0))
		} else {
			pfx = int(math.Ceil(float64(width-len(lines[row])) / 2.0))
			sfx = int(math.Floor(float64(width-len(lines[row])) / 2.0))
		}
		if surroundWithAttr {
			tbprint(lt.X, lt.Y+row, fg, bg, strings.Repeat(" ", pfx)+lines[row]+strings.Repeat(" ", sfx))
		} else {
			tbprint(lt.X+pfx, lt.Y+row, fg, bg, lines[row])
		}
	}
	return rows
}

func fill(lt Point, rb Point, fg termbox.Attribute, bg termbox.Attribute, char rune) {
	var x, y int
	for x = lt.X; x <= rb.X; x++ {
		for y = lt.Y; y <= rb.Y; y++ {
			termbox.SetCell(x, y, char, fg, bg)
		}
	}
}

func redraw() {
	termbox.Clear(coldef, coldef)
	var w, h int
	w, h = termbox.Size()
	currentScreen.draw(Point{sepWidth + padding + 1, padding}, Point{w - padding, h - padding})
	fill(Point{sepWidth, 0}, Point{sepWidth, h}, coldef, coldef, 'x')

	var scr iScreen
	var area Area
	for scr, area = range overlays {
		if area.Lt != InvalidPoint && area.Rb != InvalidPoint {
			scr.draw(area.Lt, area.Rb)
		}
	}

	termbox.Flush()
}
