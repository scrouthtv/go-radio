package main

import "testing"

import "github.com/scrouthtv/go-radio/util"
import "github.com/scrouthtv/go-radio/recorder"

func TestCompareStringSlice(t *testing.T) {
	a := []string{"asdf", "asdf", "qwertz"}
	b := []string{"asdf", "qwertz", "asdf"}
	c := []string{"asdf", "qwertz"}
	d := []string{"asdf", "huiu", "qwertz"}
	e := []string{"asdf", "asdf", "qwertz", "qwertz"}
	f := []string{"asdf", "asdf", "qwertz"}
	if util.IsStringSliceEqual(a, b) != false {
		t.Error(a, "shouldn't be equal to", b)
	}
	if util.IsStringSliceEqual(a, c) != false {
		t.Error(a, "shouldn't be equal to", c)
	}
	if util.IsStringSliceEqual(a, d) != false {
		t.Error(a, "shouldn't be equal to", d)
	}
	if util.IsStringSliceEqual(a, e) != false {
		t.Error(a, "shouldn't be equal to", e)
	}
	if util.IsStringSliceEqual(a, f) != true {
		t.Error(a, "should be equal to", e)
	}
}

func TestCompareRecordingSlice(t *testing.T) {
	a, b, c := randomRecordingsSlice(3, allOnByte), randomRecordingsSlice(5, allOnByte), randomRecordingsSlice(5, allOnByte)
	d := c
	var ca, cb, cc, cd []util.Comparable = recSliceToComparableSlice(a), recSliceToComparableSlice(b),
		recSliceToComparableSlice(c), recSliceToComparableSlice(d)

	if util.IsSliceEqual(ca, ca) != true {
		t.Error(a, "should be equal to", a)
	}
	if util.IsSliceEqual(ca, cb) != false {
		t.Error(a, "shouldn't be equal to", b)
	}
	if util.IsSliceEqual(cb, cc) != false {
		t.Error(b, "shouldn't be equal to", c)
	}
	if util.IsSliceEqual(cc, cd) != true {
		t.Error(c, "should be equal to", d)
	}
}

func recSliceToComparableSlice(rs []recorder.Recording) []util.Comparable {
	// literally the same thing again but can't reuse it
	cs := make([]util.Comparable, len(rs))
	for i, r := range rs {
		cs[i] = util.Comparable(r)
	}
	return cs
}

/*func testCompareRecording() {
	a, b := randomRecording(allOnByte), randomRecording(allOnByte)
	fmt.Println("true :", a.Equal(a))
	fmt.Println("false :", a.Equal(b))
}*/
