package main

import "testing"
import "os"
import "time"
import "bufio"
import "strings"
import "fmt"

import "github.com/scrouthtv/go-radio/recorder"
import "github.com/scrouthtv/go-radio/util"

func TestReclist(t *testing.T) {
	var normalList recorder.RecordingsList = recorder.RecordingsList{
		randomFile("*.csv"), []int{0, 1, 2, 3}, "02.01.2006 15:04:05", nil}
	var normalTest reclistTest = reclistTest{normalList, 0b11111100}
	t.Log("normal list in", normalList.Path)

	var orderedList recorder.RecordingsList = recorder.RecordingsList{
		randomFile("*.csv"), []int{2, 3, 1, 0}, "15:04:05 02.01.2006", nil}
	var orderedTest reclistTest = reclistTest{orderedList, 0b11111100}
	t.Log("ordered list in", orderedList.Path)

	var dateList recorder.RecordingsList = recorder.RecordingsList{
		randomFile("*.csv"), []int{0, 1, 2, 3}, time.RFC3339Nano, nil}
	var dateTest reclistTest = reclistTest{dateList, 0b11111111}
	t.Log("date list in", dateList.Path)

	t.Log("write randoms test")
	normalTest.writeRandoms(3, t)
	t.Log("on ordered list")
	orderedTest.writeRandoms(5, t)
	t.Log("on full date list")
	dateTest.writeRandoms(3, t)

	t.Log("delete rows test")
	normalTest.deleteLocalRow(t)
	t.Log("on file")
	normalTest.deleteFileRow(t)

	t.Log("change content test")
	orderedTest.changeFileContents(t)

	t.Log("continuity test")
	normalTest.fileContinuity(t)
}

type reclistTest struct {
	list     recorder.RecordingsList
	timemask byte
}

func (test *reclistTest) fileContinuity(t *testing.T) {
	// because idc about RAM
	var lines1, lines2, lines3, lines4, lines5, lines6 []string
	var list1, list2, list3, list4, list5, list6 []recorder.Recording

	test.writeRandoms(6, t)
	lines1 = fileLines(test.list.Path)
	list1 = test.list.Recordings

	test.list.Save()
	lines2 = fileLines(test.list.Path)
	list2 = test.list.Recordings

	test.list.Load()
	test.list.Save()
	lines3 = fileLines(test.list.Path)
	list3 = test.list.Recordings

	test.writeRandoms(6, t)
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

	if util.IsStringSliceEqual(lines1, lines2) != true {
		t.Error(lines1, "should be equal to", lines2)
	}
	if isRecSliceEqual(list1, list2) != true {
		t.Error(list1, "should be equal to", list2)
	}
	if util.IsStringSliceEqual(lines1, lines3) != true {
		t.Error(lines1, "should be equal to", lines3)
	}
	if isRecSliceEqual(list1, list3) != true {
		t.Error(list1, "should be equal to", list3)
	}
	if util.IsStringSliceEqual(lines1, lines4) != false {
		t.Error(lines1, "shouldn't be equal to", lines4)
	}
	if isRecSliceEqual(list1, list4) != false {
		t.Error(list1, "shouldn't be equal to", list4)
	}
	if util.IsStringSliceEqual(lines1, lines5) != true {
		t.Error(lines1, "should be equal to", lines5)
	}
	if isRecSliceEqual(list1, list5) != true {
		t.Error(list1, "should be equal to", list5)
	}
	if util.IsStringSliceEqual(lines4, lines6) != true {
		t.Error(lines4, "should be equal to", lines6)
	}
	if isRecSliceEqual(list4, list6) != true {
		t.Error(list4, "should be equal to", list6)
	}
}

func (test *reclistTest) changeFileContents(t *testing.T) {
	test.list.Save()
	var lines []string = fileLines(test.list.Path)
	var recordsPre string = fmt.Sprint(test.list.Recordings)

	for i, s := range lines {
		lines[i] = strings.ReplaceAll(strings.ReplaceAll(s, "q", "p"), "5", "2")
	}
	writeLines(test.list.Path, lines)
	test.list.Load()

	var recordsPost string = fmt.Sprint(test.list.Recordings)

	if recordsPre == recordsPost != false {
		t.Error(recordsPre, "shouldn't be equal to", recordsPost)
	}
	if strings.ReplaceAll(strings.ReplaceAll(recordsPre, "q", "p"), "5", "2") == recordsPost != true {
		t.Error(strings.ReplaceAll(strings.ReplaceAll(recordsPre, "q", "p"), "5", "2") == recordsPost, "should be equal to", strings.ReplaceAll(strings.ReplaceAll(recordsPre, "q", "p"), "5", "2") == recordsPost)
	}
}

func (test *reclistTest) deleteFileRow(t *testing.T) {
	test.list.Load()
	if len(test.list.Recordings) < 1 {
		t.Log("file is too short, can't delete a line")
		return
	}
	var recordsPre []recorder.Recording = test.list.Recordings

	// delete the first record:
	writeLines(test.list.Path, fileLines(test.list.Path)[1:])
	test.list.Load()

	var recordsPost []recorder.Recording = test.list.Recordings

	if isRecSliceEqual(recordsPre, recordsPost) != false {
		t.Error(recordsPre, "shouldn't be equal to", recordsPost)
	}
	if isRecSliceEqual(recordsPre[1:], recordsPost) != true {
		t.Error(recordsPre[1:], "should be equal to", recordsPost)
	}
}

func isRecSliceEqual(a []recorder.Recording, b []recorder.Recording) bool {
	return util.IsSliceEqual(recSliceToComparableSlice(a), recSliceToComparableSlice(b))
}

func (test *reclistTest) deleteLocalRow(t *testing.T) {
	test.list.Load()
	if len(test.list.Recordings) < 1 {
		t.Log("too few records, can't delete a recording")
		return
	}
	var fileLinesPre []string = fileLines(test.list.Path)

	test.list.Recordings = test.list.Recordings[1:] // remove the last line
	test.list.Save()

	var fileLinesPost []string = fileLines(test.list.Path)

	if util.IsStringSliceEqual(fileLinesPre, fileLinesPost) != false {
		t.Error(fileLinesPre, "shouldn't be equal to", fileLinesPost)
	}
	if util.IsStringSliceEqual(fileLinesPre[1:], fileLinesPost) != true {
		t.Error(fileLinesPre[1:], "should be equal to", fileLinesPost)
	}
}

func (test *reclistTest) writeRandoms(amount int, t *testing.T) {
	var errs *[]error

	var rcs []recorder.Recording = randomRecordingsSlice(amount, test.timemask)
	test.list.Recordings = rcs
	errs = test.list.Save()
	check(false, *errs...)
	// written ^

	if test.linesInRclFile() != amount {
		t.Error("Expected", amount, "lines, got", test.linesInRclFile())
	}
	test.list.Recordings = nil
	errs = test.list.Load()
	check(false, *errs...)
	if amount != len(test.list.Recordings) {
		t.Error("Expected", amount, "records, got", len(test.list.Recordings))
	}
	if isRecSliceEqual(rcs, test.list.Recordings) != true {
		t.Error(rcs, "should be equal to", test.list.Recordings)
	}
}

func (test *reclistTest) linesInRclFile() int {
	var f *os.File
	f, _ = os.Open(test.list.Path)
	var lines int = 0
	var rdr *bufio.Scanner = bufio.NewScanner(f)
	for rdr.Scan() {
		lines++
	}
	f.Close()
	return lines
}

func (test *reclistTest) dumpRecordings(t *testing.T) {
	var i int
	var rc recorder.Recording

	for i, rc = range test.list.Recordings {
		t.Log(i, rc.String())
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
