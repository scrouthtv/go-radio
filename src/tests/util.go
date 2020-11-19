package main

import "math/rand"
import "strings"
import "net/url"
import "time"
import "fmt"

func main() {
	zerotime, _ := time.Parse("", "")
	fmt.Println("true :", zerotime.Equal(randomTime(0)))
	testReclist()
	//testCompare()
}

func randomInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

const defChars string = "abcdefgjhijklmnopqrstuvwxyz"

func randomRune(runes []rune) rune {
	return runes[rand.Intn(len(runes))]
}

func randomBool() bool {
	return rand.Intn(2) == 0
}

func randomTime(mask byte) time.Time {
	var year, day, hour, minute, sec, nsec = 0, 1, 0, 0, 0, 0
	var month time.Month = time.January
	var zone *time.Location = time.UTC

	if mask&0b10000000 > 0 {
		year = randomInt(2000, 2030)
	}
	if mask&0b01000000 > 0 {
		month = time.Month(randomInt(1, 12))
	}
	if mask&0b00100000 > 0 {
		day = rand.Intn(29)
	}
	if mask&0b00010000 > 0 {
		hour = rand.Intn(24)
	}
	if mask&0b00001000 > 0 {
		minute = rand.Intn(60)
	}
	if mask&0b00000100 > 0 {
		sec = rand.Intn(60)
	}
	if mask&0b0000010 > 0 {
		nsec = rand.Intn(1000000000)
	}
	if mask&0b0000001 > 0 {
		zone, _ = time.LoadLocation(randomElement(zones))
	}
	return time.Date(year, month, day, hour, minute, sec, nsec, zone)
}

func randomElement(arr []string) string {
	return arr[rand.Intn(len(arr))]
}

// proto should be one of http, https, ...
// might return nil
func randomURL(proto string) *url.URL {
	var domain, tld string = randomString(5, 12), randomString(2, 3)
	var sfx string

	var folders int = rand.Intn(5)
	var i int
	for i = 0; i < folders; i++ {
		sfx += randomString(2, 16)
		if i < folders-1 {
			sfx += "/"
		} else if randomBool() {
			sfx += "." + randomString(2, 4)
		}
	}

	var myurl *url.URL
	myurl, _ = url.Parse(proto + "://" + domain + "." + tld + "/" + sfx)
	return myurl
}

func randomString(minlen int, maxlen int) string {
	var strlen int = randomInt(minlen, maxlen)
	var i int
	var str strings.Builder
	for i = 0; i < strlen; i++ {
		str.WriteRune(randomRune([]rune(defChars)))
	}
	return str.String()
}

func check(doExit bool, errs ...error) {
	var i int
	var err error
	var failed bool = false
	for i, err = range errs {
		if err != nil {
			fmt.Println(i, err)
			failed = true
		}
	}
	if failed && doExit {
		panic("Error.")
	}
}
