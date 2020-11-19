package main

import "fmt"
import "io/ioutil"
import "os"
import "time"
import "bufio"

import "github.com/scrouthtv/go-radio/recorder"
import "github.com/scrouthtv/go-radio/util"

func testReclist() {
	var file *os.File
	var err error
	file, err = ioutil.TempFile("", "*.csv")
	file.Close()
	check(true, err)
	fmt.Println("normal list in", file.Name())
	var normalList recorder.RecordingsList = recorder.RecordingsList{
		file.Name(), []int{0, 1, 2, 3}, "02.01.2006 15:04:05", nil}
	var normalTest reclistTest = reclistTest{normalList, 0b11111100}

	file, err = ioutil.TempFile("", "*.csv")
	file.Close()
	check(true, err)
	fmt.Println("ordered list in", file.Name())
	var orderedList recorder.RecordingsList = recorder.RecordingsList{
		file.Name(), []int{2, 3, 1, 0}, "15:04:05 02.01.2006", nil}
	var orderedTest reclistTest = reclistTest{orderedList, 0b11111100}

	file, err = ioutil.TempFile("", "*.csv")
	file.Close()
	check(true, err)
	fmt.Println("long date list in", file.Name())
	var dateList recorder.RecordingsList = recorder.RecordingsList{
		file.Name(), []int{0, 1, 2, 3}, time.RFC3339Nano, nil}
	var dateTest reclistTest = reclistTest{dateList, 0b11111111}

	normalTest.writeRandoms(0)
	orderedTest.writeRandoms(2)
	dateTest.writeRandoms(0)
}

type reclistTest struct {
	list     recorder.RecordingsList
	timemask byte
}

func (test reclistTest) writeRandoms(amount int) {
	var errs *[]error

	var rcs []recorder.Recording = randomRecordingsSlice(amount, test.timemask)
	test.list.Recordings = rcs
	errs = test.list.Save()
	check(false, *errs...)
	// written ^

	fmt.Print(amount, " : ")
	test.linesInRclFile()
	test.list.Recordings = nil
	errs = test.list.Load()
	check(false, *errs...)
	fmt.Println(amount, ":", len(test.list.Recordings))
	fmt.Println("true :", util.IsSliceEqual(toComparableSlice(rcs), toComparableSlice(test.list.Recordings)))
	if !util.IsSliceEqual(toComparableSlice(rcs), toComparableSlice(test.list.Recordings)) {
		/*fmt.Println(rcs)
		fmt.Println(test.list.Recordings)*/
	}
}

func (test reclistTest) linesInRclFile() {
	var f *os.File
	f, _ = os.Open(test.list.Path)
	var lines int = 0
	var rdr *bufio.Scanner = bufio.NewScanner(f)
	for rdr.Scan() {
		lines++
	}
	fmt.Println(lines)
	f.Close()
}

func (test reclistTest) dumpRecordings() {
	var i int
	var rc recorder.Recording

	for i, rc = range test.list.Recordings {
		fmt.Println(i, rc.String())
	}
}

func randomRecording(timemask byte) recorder.Recording {
	var r recorder.Recording
	r.Enabled = randomBool()
	r.Stream = randomURL("https")

	var startTime, endTime time.Time
	startTime, endTime = randomTime(timemask), randomTime(timemask)
	r.Start = &startTime
	r.End = &endTime

	return r
}

func randomRecordingsSlice(amount int, timemask byte) []recorder.Recording {
	var rcs []recorder.Recording

	var i int
	for i = 0; i < amount; i++ {
		rcs = append(rcs, randomRecording(timemask))
	}

	return rcs
}
