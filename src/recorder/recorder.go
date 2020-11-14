package recorder

type Recorder interface {
	Record(stream string, filepath string, seconds int) error
}
