package main

import "fmt"

import "github.com/scrouthtv/go-radio/util"
import "github.com/scrouthtv/go-radio/recorder"

const allOnByte byte = 0b11111111

func testCompare() {
	testCompareRecording()
	testCompareRecordingSlice()
}

func testCompareRecordingSlice() {
	a, b, c := randomRecordingsSlice(3, allOnByte), randomRecordingsSlice(5, allOnByte), randomRecordingsSlice(5, allOnByte)
	d := c
	var ca, cb, cc, cd []util.Comparable = toComparableSlice(a), toComparableSlice(b),
		toComparableSlice(c), toComparableSlice(d)

	fmt.Println("true :", util.IsSliceEqual(ca, ca))
	fmt.Println("false :", util.IsSliceEqual(ca, cb))
	fmt.Println("false :", util.IsSliceEqual(cb, cc))
	fmt.Println("true :", util.IsSliceEqual(cc, cd))
}

func toComparableSlice(rs []recorder.Recording) []util.Comparable {
	cs := make([]util.Comparable, len(rs))
	for i, r := range rs {
		cs[i] = util.Comparable(r)
	}
	return cs
}

func testCompareRecording() {
	a, b := randomRecording(allOnByte), randomRecording(allOnByte)
	fmt.Println("true :", a.Equal(a))
	fmt.Println("false :", a.Equal(b))
}
