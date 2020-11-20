package recorder

import "io"
import "encoding/csv"
import "net/url"
import "time"
import "os"
import "fmt"

import "github.com/scrouthtv/go-radio/util"

type Recording struct {
	Enabled bool
	Stream  *url.URL
	Start   *time.Time
	End     *time.Time
}

func (this Recording) Equal(other util.Comparable) bool {
	var rec Recording
	var ok bool
	rec, ok = other.(Recording)
	if !ok {
		return false
	} else if other == nil {
		return false
	} else if this.Enabled != rec.Enabled {
		return false
	} else if this.Stream.String() != rec.Stream.String() {
		return false
	} else {
		return this.Start.Equal(*rec.Start) && this.End.Equal(*rec.End)
	}
}

type RecordingsList struct {
	Path string
	// in which order the fields enabled, stream, start, end appear in the file:
	FieldOrder []int
	TimeFormat string
	Recordings []Recording
	// was *[], []* makes more sense but is stupid: https://philpearl.github.io/post/bad_go_slice_of_pointers/
}

func NewRecordingsList(path string, fieldOrder []int, timeFormat string) (*RecordingsList, error) {
	var err error
	path, err = util.ExpandPath(path)
	if err != nil {
		return nil, err
	}
	var list RecordingsList = RecordingsList{path, fieldOrder, timeFormat, nil}
	return &list, nil
}

func (rec Recording) String() string {
	return fmt.Sprintf("Enabled: %s, via %s,\n%s - %s", formatBool(&rec.Enabled),
		rec.Stream.String(), rec.Start.String(), rec.End.String())
}

var DefaultRecordingsList, _ = NewRecordingsList(
	"~/.config/go-radio/recordings.csv",
	[]int{0, 1, 2, 3},
	"02.01.2006 15:04:05",
)

func (list *RecordingsList) Load() *[]error {
	var err error
	var errs []error
	list.Recordings = nil // empty the array in memory
	var f *os.File
	f, err = os.OpenFile(list.Path, os.O_RDONLY, 0644)
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
			// last record is nil & io.EOF, return:
			f.Close()
			return &errs
		} else if err != nil {
			errs = append(errs, err)
			continue
		}
		rcd.Enabled = parseBool(&rc[list.FieldOrder[0]])
		rcd.Stream, err = url.Parse(rc[list.FieldOrder[1]])
		if err != nil {
			errs = append(errs, err)
			continue
		}

		var startTime, endTime time.Time
		startTime, err = time.Parse(list.TimeFormat, rc[list.FieldOrder[2]])
		if err != nil {
			errs = append(errs, err)
			continue
		}
		rcd.Start = &startTime
		endTime, err = time.Parse(list.TimeFormat, rc[list.FieldOrder[3]])
		if err != nil {
			errs = append(errs, err)
			continue
		}
		rcd.End = &endTime

		list.Recordings = append(list.Recordings, rcd)
	}
	return &errs // we're never coming here anyways
}

// Writes a file with the specified recordings
// Ignores fields that are not specified in the FieldOrder
func (list *RecordingsList) Save() *[]error {
	var err error
	var errs []error
	var f *os.File
	f, err = os.OpenFile(list.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		errs = append(errs, err)
		return &errs
	}

	var wr *csv.Writer = csv.NewWriter(f)
	var rc []string = make([]string, len(list.FieldOrder))
	var rcd Recording
	for _, rcd = range list.Recordings {
		rc[list.FieldOrder[0]] = formatBool(&rcd.Enabled)
		rc[list.FieldOrder[1]] = rcd.Stream.String()
		rc[list.FieldOrder[2]] = rcd.Start.Format(list.TimeFormat)
		rc[list.FieldOrder[3]] = rcd.End.Format(list.TimeFormat)

		err = wr.Write(rc)
		if err != nil {
			errs = append(errs, err)
		}
		err = wr.Error()
		if err != nil {
			errs = append(errs, err)
		}
	}

	wr.Flush()
	//f.Close()
	return &errs
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
