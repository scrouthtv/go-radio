package util

type Comparable interface {
	Equal(Comparable) bool
}

func IsStringSliceEqual(arr1 []string, arr2 []string) bool {
	if len(arr1) == len(arr2) {
		var i int
		var c1, c2 string
		for i, c1 = range arr1 {
			c2 = arr2[i]
			if c1 != c2 {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func IsSliceEqual(arr1 []Comparable, arr2 []Comparable) bool {
	if len(arr1) == len(arr2) {
		var i int
		var c1, c2 Comparable
		for i, c1 = range arr1 {
			c2 = arr2[i]
			if !c1.Equal(c2) {
				return false
			}
		}
		return true
	} else {
		return false
	}
}
