package tui

import "strings"
import "sort"
import "time"

import "github.com/mitchellh/go-wordwrap" // pretty schtoopid to use a whole lib for this single purpose

import "github.com/scrouthtv/go-radio/stations"

func softwrap(msg string, width int) []string {
	return strings.Split(wordwrap.WrapString(msg, uint(width)), "\n")
	var i int
	var line string
	var lines []string
	var word string
	for i, word = range strings.Split(msg, " ") {
		if len(line)+len(word)+1 < width { // can simply append
			if i == 0 {
				line = word
			} else {
				line += " " + word
			}
		} else if len(word) > width { // word is longer than a full line
			var stripped int
			for len(word) > width {
				stripped = width - len(line) - 1
				if stripped > len(word)-1 {
					stripped = len(word) - 1
				}
				if stripped < 0 {
					// TODO this is super bugged
					return lines
					log("")
					log("")
					log("")
					log("")
					log("")
					log("")
					log("")
					log("")
					log("")
					log("")
					log("")
					log(lines)
					log("current word: ", word)
					log("want to strip ", stripped)
					log("current line: ", line)
				}
				line += " " + word
				lines = append(lines, line)
				line = ""
				word = word[stripped:len(word)]
			}
			line = " " + word
		} else { // have to make a new line
			lines = append(lines, line)
			line = word
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

func minutesToHM(minutes int) (int, int) {
	return minutes / 60, minutes % 60
}

func SortEventsByStart(events []stations.Event) []stations.Event {
	sort.SliceStable(events, func(i1 int, i2 int) bool {
		return events[i2].Start.After(events[i1].Start)
	})
	return events
}

func shrink(str string, length int) string {
	if length <= 0 {
		return ""
	} else if len(str) > 3 && length <= 3 {
		return "..."
	}
	if len(str) > length {
		return str[0:length-3] + "..."
	}
	return str
}

func timeAsMinutes(time time.Time) int {
	return time.Hour()*60 + time.Minute()
}

func (p Point) Add(x int, y int) Point {
	return Point{p.X + x, p.Y + y}
}

func Contains(arr []rune, x rune) bool {
	var a interface{}
	for _, a = range arr {
		if a == x {
			return true
		}
	}
	return false
}

// for black value calculation
func IsFullRune(r rune) bool {
	var halfRunes []rune = []rune{'.', ' ', '-', ':', ','}
	return !Contains(halfRunes, r)
}
