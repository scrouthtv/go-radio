package stations

import "time"

//														 YMDhmsnloc
const combineMaskYMD uint8 = 0b11100000

func combine(a time.Time, b time.Time, mask uint8) time.Time {
	var year, day int
	var month time.Month
	var hour, minute, second, nanosecond int
	var loc *time.Location

	if mask&0b10000000 > 0 {
		year = a.Year()
	} else {
		year = b.Year()
	}
	if mask&0b01000000 > 0 {
		month = a.Month()
	} else {
		month = b.Month()
	}
	if mask&0b00100000 > 0 {
		day = a.Day()
	} else {
		day = b.Day()
	}
	if mask&0b00010000 > 0 {
		hour = a.Hour()
	} else {
		hour = b.Hour()
	}
	if mask&0b00001000 > 0 {
		minute = a.Minute()
	} else {
		minute = b.Minute()
	}
	if mask&0b00000100 > 0 {
		second = a.Second()
	} else {
		second = b.Second()
	}
	if mask&0b00000010 > 0 {
		nanosecond = a.Nanosecond()
	} else {
		nanosecond = b.Nanosecond()
	}
	if mask&0b00000001 > 0 {
		loc = a.Location()
	} else {
		loc = b.Location()
	}

	return time.Date(year, month, day, hour, minute, second, nanosecond, loc)
}
