package main

import "io"
import "encoding/csv"
import "url"

type Recording struct {
	Enabled bool
	Stream  *url.URL
	Start   *time.Time
	End     *time.Time
}

type RecordingsList struct {
	// FULL path, no ~ or variables
	Path string
	// in which order the fields enabled, stream, start, end appear in the file:
	FieldOrder []int
	TimeFormat string
	Recordings *[]Recording
}

const DefaultRecordingsList = RecordingsList{
	"/home/lenni/.config/go-radio/recordings.csv",
	[]int{0, 1, 2, 3},
	"02.01.2006 15:04:05",
}

func (list RecordingsList) Load() *[]error {
	var errs []error
	list.Recordings = nil // empty the array in memory
	var f *os.File
	f, err = os.OpenFile(list.path, os.O_RDONLY, 0644)
	if err != nil {
		errs = append(errs, err)
		return &errs
	}

	var rdr *csv.Reader = csv.NewReader(f)
	var rc []string // a single record
	var rcd Recording
	for {
		rc, err = rdr.Read()
		if err == io.EOF {
			// return:
			return &errs
		} else if err != nil {
			errs = append(errs, err)
			continue
		}
		rcd.Enabled = parseBool(rc[list.FieldOrder[0]])
		rcd.Stream, err = url.Parse(rc[list.FieldOrder[1]])
		if err != nil {
			errs = append(errs, err)
			continue
		}
		rcd.Start, err = time.Parse(list.TimeFormat, rc[list.FieldOrder[2]])
		if err != nil {
			errs = append(errs, err)
			continue
		}
		rcd.Time, err = time.Parse(list.TimeFormat, rc[list.FieldOrder[2]])
		if err != nil {
			errs = append(errs, err)
			continue
		}

		list.Recordings = append(list.Recordings, rcd)
	}
}

// if FieldOrder is e. g. [ 0, 2, 3, 1 ], that means that the actual file itself
// looks like this: enabled,start,end,stream .
// Now if I want to save the start field for example, I need to know at which
// position it is (in this example 1). This function is basically the inverse of
// the FieldOrder list
func (list RecordingsList) firstPositionFor(searchField int) int {
	var pos, field int
	for pos, field = range list.FieldOrder {
		if field == searchField {
			return pos
		}
	}
	return -1
}

// Writes a file with the specified recordings
// Ignores fields that are not specified in the FieldOrder
func (list RecordingsList) Save() *[]error {
	var errs []error
	var f *os.File
	f, err = os.OpenFile(list.Path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		errs = append(errs, err)
		return &errs
	}

	var wr *csv.Writer = csv.NewWriter(f)
	var rc []string = make([]string, len(list.FieldOrder))
	var rcd Recording
	var pos int
	for _, rcd = range &list.Recordings {
		pos = list.firstPositionFor(0)
		if pos > 0 {
			rc[pos] = formatBool(rcd.Enabled)
		}
		pos = list.firstPositionFor(1)
		if pos > 0 {
			rc[pos] = rcd.Stream.String()
		}
		pos = list.firstPositionFor(2)
		if pos > 0 {
			rc[pos] = rcd.Start.Format(list.TimeFormat)
		}
		pos = list.firstPositionFor(3)
		if pos > 0 {
			rc[pos] = rcd.End.Format(list.TimeFormat)
		}
	}
}

func formatBool(b *bool) string {
	if *b {
		return "true"
	} else {
		return "false"
	}
}

func parseBool(str *string) bool {
	return *str == "1" || *str == "on" || *str == "true" || *str == "yes"
}
