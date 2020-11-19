package main

import "fmt"
import "os"
import "time"
import "bufio"
import "strings"

import "github.com/scrouthtv/go-radio/recorder"
import "github.com/scrouthtv/go-radio/util"

func testReclist() {
	var normalList recorder.RecordingsList = recorder.RecordingsList{
		randomFile("*.csv"), []int{0, 1, 2, 3}, "02.01.2006 15:04:05", nil}
	var normalTest reclistTest = reclistTest{normalList, 0b11111100}
	fmt.Println("normal list in", normalList.Path)

	var orderedList recorder.RecordingsList = recorder.RecordingsList{
		randomFile("*.csv"), []int{2, 3, 1, 0}, "15:04:05 02.01.2006", nil}
	var orderedTest reclistTest = reclistTest{orderedList, 0b11111100}
	fmt.Println("ordered list in", orderedList.Path)

	var dateList recorder.RecordingsList = recorder.RecordingsList{
		randomFile("*.csv"), []int{0, 1, 2, 3}, time.RFC3339Nano, nil}
	var dateTest reclistTest = reclistTest{dateList, 0b11111111}
	fmt.Println("date list in", dateList.Path)

	fmt.Println("write randoms test")
	normalTest.writeRandoms(3)
	fmt.Println("on ordered list")
	orderedTest.writeRandoms(5)
	fmt.Println("on full date list")
	dateTest.writeRandoms(3)

	fmt.Println("delete rows test")
	normalTest.deleteLocalRow()
	fmt.Println("on file")
	normalTest.deleteFileRow()

	fmt.Println("change content test")
	orderedTest.changeFileContents()

	fmt.Println("continuity test")
	normalTest.fileContinuity()
}

type reclistTest struct {
	list     recorder.RecordingsList
	timemask byte
}

func (test *reclistTest) fileContinuity() {
	// because idc about RAM
	var lines1, lines2, lines3, lines4, lines5, lines6 []string
	var list1, list2, list3, list4, list5, list6 []recorder.Recording

	test.writeRandoms(6)
	lines1 = fileLines(test.list.Path)
	list1 = test.list.Recordings

	test.list.Save()
	lines2 = fileLines(test.list.Path)
	list2 = test.list.Recordings

	test.list.Load()
	test.list.Save()
	lines3 = fileLines(test.list.Path)
	list3 = test.list.Recordings

	test.writeRandoms(6)
	lines4 = fileLines(test.list.Path)
	list4 = test.list.Recordings

	writeLines(test.list.Path, lines2)
	test.list.Load()
	lines5 = fileLines(test.list.Path)
	list5 = test.list.Recordings

	writeLines(test.list.Path, lines4)
	test.list.Load()
	lines6 = fileLines(test.list.Path)
	list6 = test.list.Recordings

	fmt.Println("test finished")

	fmt.Println("true :", util.IsStringSliceEqual(lines1, lines2))
	fmt.Println("true :", isRecSliceEqual(list1, list2))
	fmt.Println("true :", util.IsStringSliceEqual(lines1, lines3))
	fmt.Println("true :", isRecSliceEqual(list1, list3))
	fmt.Println("false :", util.IsStringSliceEqual(lines1, lines4))
	fmt.Println("false :", isRecSliceEqual(list1, list4))
	fmt.Println("true :", util.IsStringSliceEqual(lines1, lines5))
	fmt.Println("true :", isRecSliceEqual(list1, list5))
	fmt.Println("true :", util.IsStringSliceEqual(lines4, lines6))
	fmt.Println("true :", isRecSliceEqual(list4, list6))
}

func (test *reclistTest) changeFileContents() {
	test.list.Save()
	var lines []string = fileLines(test.list.Path)
	var recordsPre string = fmt.Sprint(test.list.Recordings)

	for i, s := range lines {
		lines[i] = strings.ReplaceAll(strings.ReplaceAll(s, "q", "p"), "5", "2")
	}
	writeLines(test.list.Path, lines)
	test.list.Load()

	var recordsPost string = fmt.Sprint(test.list.Recordings)

	fmt.Println("false :", recordsPre == recordsPost)
	fmt.Println("true :", strings.ReplaceAll(strings.ReplaceAll(recordsPre, "q", "p"), "5", "2") == recordsPost)
}

func (test *reclistTest) deleteFileRow() {
	test.list.Load()
	if len(test.list.Recordings) < 1 {
		fmt.Println("file is too short, can't delete a line")
		return
	}
	var recordsPre []recorder.Recording = test.list.Recordings

	// delete the first record:
	writeLines(test.list.Path, fileLines(test.list.Path)[1:])
	test.list.Load()

	var recordsPost []recorder.Recording = test.list.Recordings

	fmt.Println("false :", isRecSliceEqual(recordsPre, recordsPost))
	fmt.Println("true :", isRecSliceEqual(recordsPre[1:], recordsPost))
}

func isRecSliceEqual(a []recorder.Recording, b []recorder.Recording) bool {
	return util.IsSliceEqual(recSliceToComparableSlice(a), recSliceToComparableSlice(b))
}

func (test *reclistTest) deleteLocalRow() {
	test.list.Load()
	if len(test.list.Recordings) < 1 {
		fmt.Println("too few records, can't delete a recording")
		return
	}
	var fileLinesPre []string = fileLines(test.list.Path)

	test.list.Recordings = test.list.Recordings[1:] // remove the last line
	test.list.Save()

	var fileLinesPost []string = fileLines(test.list.Path)

	fmt.Println("false :", util.IsStringSliceEqual(fileLinesPre, fileLinesPost))
	fmt.Println("true :", util.IsStringSliceEqual(fileLinesPre[1:], fileLinesPost))
}

func (test *reclistTest) writeRandoms(amount int) {
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
	fmt.Println("true :", isRecSliceEqual(rcs, test.list.Recordings))
}

func (test *reclistTest) linesInRclFile() {
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

func (test *reclistTest) dumpRecordings() {
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
