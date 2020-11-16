package main

//import "github.com/scrouthtv/go-radio/tui"
import "github.com/scrouthtv/go-radio/recorder"
import "github.com/scrouthtv/go-radio/stations"

import "time"
import "fmt"
import "os"

var r recorder.Recorder

func main() {
	//go tui.DrawThread()
	r = recorder.DownloadRecorder{}
	if len(os.Args) < 2 {
		fmt.Println("not enough args")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "download-now":
		var stream string = "https://swr-edge-2032-dus-lg-cdn.cast.addradio.de/swr/swr/raka03/mp3/128/stream.mp3"
		var out string = "/tmp/stream.mp3"
		var end time.Time = time.Now().Add(time.Second * 10)
		var err error = r.Record(stream, out, end)
		if err == nil {
			fmt.Println("no error reported")
		} else {
			fmt.Println(err)
		}
	case "dlf-programm":
		events, err := stations.Deutschlandfunk.DailyProgram(time.Now())
		if err != nil {
			fmt.Println("err:")
			fmt.Println(err)
			os.Exit(1)
		}
		var i int
		var ev stations.Event
		for i, ev = range events {
			fmt.Printf("%3d: %s from %02d:%02d to %02d:%02dh, Category %s", i, ev.Name, ev.Start.Hour(),
				ev.Start.Minute(), ev.End.Hour(), ev.End.Minute(), ev.Category)
		}
	default:
		fmt.Println("Unknown option", os.Args[1])
	}
}
