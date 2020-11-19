package main

import "math/rand"
import "strings"
import "net/url"
import "time"
import "fmt"

func main() {
	testReclist()
	//testCsv()
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

func randomTime() time.Time {
	var year int = randomInt(2000, 2030)
	var month time.Month = time.Month(randomInt(1, 12))
	var day int = rand.Intn(29)
	var hour, minute, sec, nsec = rand.Intn(24), rand.Intn(60), rand.Intn(60), rand.Intn(1000000000)
	var zone *time.Location
	zone, _ = time.LoadLocation(randomElement(zones))
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
