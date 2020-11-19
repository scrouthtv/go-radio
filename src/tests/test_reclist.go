package main

import "fmt"
import "io/ioutil"
import "os"
import "time"
import "bufio"

import "github.com/scrouthtv/go-radio/recorder"

func testReclist() {
	var file *os.File
	var err error
	file, err = ioutil.TempFile("", "*.csv")
	file.Close()
	check(true, err)

	fmt.Println("normal list in", file.Name())
	var normalList recorder.RecordingsList = recorder.RecordingsList{
		file.Name(), []int{0, 1, 2, 3}, "02.01.2006 15:04:05", nil,
	}

	var errs *[]error

	normalList.Recordings = randomRecordingsSlice(3)
	fmt.Println("recordings in normal list:")
	fmt.Println(len(normalList.Recordings))
	//dumpRecordings(normalList.Recordings)
	fmt.Println("--")
	errs = normalList.Save()
	check(false, *errs...)
	statRclFile(file.Name())

	file, err = ioutil.TempFile("", "*.csv")
	file.Close()
	check(true, err)

	fmt.Println("ordered list in", file.Name())
	var orderedList recorder.RecordingsList = recorder.RecordingsList{
		file.Name(), []int{2, 3, 1, 0}, "15:04:05 02.01.2006", nil,
	}

	orderedList.Recordings = randomRecordingsSlice(2)
	fmt.Println("recordings in ordered list:")
	//dumpRecordings(orderedList.Recordings)
	fmt.Println(len(orderedList.Recordings))
	fmt.Println("--")
	errs = orderedList.Save()
	check(false, *errs...)
	statRclFile(file.Name())
}

func statRclFile(file string) {
	var f *os.File
	f, _ = os.Open(file)
	var lines int = 0
	var rdr *bufio.Scanner = bufio.NewScanner(f)
	for rdr.Scan() {
		lines++
	}
	fmt.Println(lines, "lines")
	f.Close()
}

func randomRecording() recorder.Recording {
	var r recorder.Recording
	r.Enabled = randomBool()
	r.Stream = randomURL("https")

	var startTime, endTime time.Time
	startTime, endTime = randomTime(), randomTime()
	r.Start = &startTime
	r.End = &endTime

	return r
}

func dumpRecordings(rcs []recorder.Recording) {
	var i int
	var rc recorder.Recording

	for i, rc = range rcs {
		fmt.Println(i, rc.String())
	}
}

func randomRecordingsSlice(amount int) []recorder.Recording {
	var rcs []recorder.Recording

	var i int
	for i = 0; i < amount; i++ {
		rcs = append(rcs, randomRecording())
	}

	return rcs
}
