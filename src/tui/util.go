package tui

import "strings"

func softwrap(msg string, width int) []string {
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
