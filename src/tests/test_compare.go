package main

import "fmt"

import "github.com/scrouthtv/go-radio/util"
import "github.com/scrouthtv/go-radio/recorder"

func testCompare() {
	testCompareStringSlice()
	testCompareRecording()
	testCompareRecordingSlice()
}

func testCompareStringSlice() {
	a := []string{"asdf", "asdf", "qwertz"}
	b := []string{"asdf", "qwertz", "asdf"}
	c := []string{"asdf", "qwertz"}
	d := []string{"asdf", "huiu", "qwertz"}
	e := []string{"asdf", "asdf", "qwertz", "qwertz"}
	f := []string{"asdf", "asdf", "qwertz"}
	fmt.Println("false :", util.IsStringSliceEqual(a, b))
	fmt.Println("false :", util.IsStringSliceEqual(a, c))
	fmt.Println("false :", util.IsStringSliceEqual(a, d))
	fmt.Println("false :", util.IsStringSliceEqual(a, e))
	fmt.Println("true :", util.IsStringSliceEqual(a, f))
}

func testCompareRecordingSlice() {
	a, b, c := randomRecordingsSlice(3, allOnByte), randomRecordingsSlice(5, allOnByte), randomRecordingsSlice(5, allOnByte)
	d := c
	var ca, cb, cc, cd []util.Comparable = recSliceToComparableSlice(a), recSliceToComparableSlice(b),
		recSliceToComparableSlice(c), recSliceToComparableSlice(d)

	fmt.Println("true :", util.IsSliceEqual(ca, ca))
	fmt.Println("false :", util.IsSliceEqual(ca, cb))
	fmt.Println("false :", util.IsSliceEqual(cb, cc))
	fmt.Println("true :", util.IsSliceEqual(cc, cd))
}

func recSliceToComparableSlice(rs []recorder.Recording) []util.Comparable {
	// literally the same thing again but can't reuse it
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
