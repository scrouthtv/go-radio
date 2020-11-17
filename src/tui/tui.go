package tui

import "fmt"
import "time"

import "github.com/nsf/termbox-go"

import "github.com/scrouthtv/go-radio/stations"

type iScreen interface {
	title() string
	draw(lt Point, rb Point)
	event(termbox.Event) bool
}

type Point struct {
	X int
	Y int
}

var Running bool = false

var currentScreen iScreen
var focusScreen iScreen

const coldef termbox.Attribute = termbox.ColorDefault

var sepWidth int

const padding int = 1

func TuiLoop() error {
	Running = true
	var err error
	err = termbox.Init()
	if err != nil {
		return err
	}

	sepWidth = 8
	currentScreen = &timelineScreen{stations.Deutschlandfunk, time.Now()}
	focusScreen = &timelineScreen{stations.Deutschlandfunk, time.Now()}

	// TODO mouse support?

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
		log("draw loop")
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
	termbox.Flush()
}

var InvalidPoint Point = Point{-1, -1}
