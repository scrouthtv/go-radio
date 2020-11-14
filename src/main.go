package main

//import "github.com/scrouthtv/go-radio/tui"
import "github.com/scrouthtv/go-radio/recorder"

var r recorder.Recorder

func main() {
	//go tui.DrawThread()
	r = recorder.DownloadRecorder{}
	switch os.Args[0] {
	case "download-now":
		var stream string = "https://swr-edge-2032-dus-lg-cdn.cast.addradio.de/swr/swr/raka03/mp3/128/stream.mp3"
		var out string = "/tmp/stream.mp3"
		var err error = r.Record(stream, out, 15)
	}
}
