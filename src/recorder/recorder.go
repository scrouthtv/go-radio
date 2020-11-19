package recorder

import "time"

type Recorder interface {
	Record(stream string, filepath string, end time.Time) error
}
