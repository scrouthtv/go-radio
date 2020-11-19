package main

import "github.com/scrouthtv/go-radio/tui"
import "github.com/scrouthtv/go-radio/recorder"
import "github.com/scrouthtv/go-radio/stations"

import "time"
import "fmt"
import "os"
import "strings"

var r recorder.Recorder

func main() {
	//go tui.DrawThread()

	if true == false {
		strings.Split("", "")
	}

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
	case "tui":
		go tui.TuiLoop()
		time.Sleep(1 * time.Second)
		for tui.Running {
			time.Sleep(1 * time.Second)
		}
		os.Exit(0)
	case "dlf-programm":
		events, err := stations.Deutschlandfunk.DailyProgram(time.Now())
		if err != nil {
			fmt.Println("err:")
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Println("got", len(events), "events:")
			tui.SortEventsByStart(events)
			var i int
			var ev stations.Event
			for i, ev = range events {
				fmt.Printf("%3d: %02d:%02d - %02d:%02dh \"%s\", more via %s, Category %s (%s)\n", i, ev.Start.Hour(),
					ev.Start.Minute(), ev.End.Hour(), ev.End.Minute(), ev.Name, ev.Url.String(), ev.Category, ev.CatUrl.String())
				if ev.Info != "" {
					var line string
					for _, line = range strings.Split(ev.Info, "\n") {
						fmt.Println("  >", line)
					}
				}
			}
		}
	/*case "splittest":
			const msg string = `Weniger Einnahmen, mehr Zuspruch?
	Die Kirchen in der Corona-Krise
	Am Mikrofon: Michael Roehl

	Gäste:
	Dr. Johann Hinrich Claussen, Kulturbeauftragter des Rates der EKD
	Rainer Schießler, Pfarrer, Pfarrei St. Maximilian München, München
	Corinna Zisselsberger, Pfarrerin, Evangelische Kirchengemeinde St. Petri-St. Marien, Berlin
	Dr. Christiane Florin, Deutschlandfunk-Redakteurin für Religion und Gesellschaft und freie Autorin

	Hörertel.: 00800 - 4464 4464
	laenderzeit@deutschlandfunk.deDie großen deutschen Kirchen leiden unter einem eklatanten Mitgliederschwund. Weit über zweihunderttausend Katholiken und ebenso viele Protestanten erklärten im letzten Jahr ihren Kirchenaustritt. Und der Trend hält an. Bis zum Jahr 2060 könnte sich die Zahl der Gläubigen sogar halbieren, schätzen Wissenschaftler. Ein riesiges Problem, das die Kirchen mittel- und langfristig in ihrer Existenz bedroht. Und nun sorgt auch noch die Corona-Pandemie für große Herausforderungen. Gottesdienste dürfen nur noch unter Sicherheitsauflagen und mit stark verminderter Besucherzahl durchgeführt werden. Oft wurden sie komplett durch Videoübertragungen aus den Gotteshäusern ersetzt. Gemeindearbeit und Seelsorge sind stark beeinträchtigt. Doch gerade in Krisenzeiten ist geistlicher Zuspruch gefragt und viele Gläubige fühlen sich vernachlässigt und im Stich gelassen. Mit innovativen Ideen versuchen manche Kirchenleute den Kontakt zu den Menschen in ihrer Gemeinde dennoch aufrecht zu erhalten.`
			const width int = 10
			var lines []string = tui.Softwrap*/
	default:
		fmt.Println("Unknown option", os.Args[1])
	}
}
