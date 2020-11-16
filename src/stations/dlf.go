package stations

import "net/http"
import "time"
import "fmt"
import "strings"

import "github.com/PuerkitoBio/goquery"

type Dlf struct {
	name       string
	url        string
	programurl string
}

var Deutschlandfunk Dlf = Dlf{"Deutschlandfunk", "https://st01.sslstream.dlf.de/dlf/01/high/aac/stream.aac", "https://www.deutschlandfunk.de/programmvorschau.281.de.html"}
var DlfKultur Dlf = Dlf{"Deutschlandfunk", "https://st02.sslstream.dlf.de/dlf/02/high/aac/stream.aac", "https://www.deutschlandfunkkultur.de/programmvorschau.282.de.html"}

func (dlf Dlf) GetName() (string, error) {
	return dlf.name, nil
}

func (dlf Dlf) GetURL() (string, error) {
	return dlf.url, nil
}

func (dlf Dlf) Program() ([]Event, error) {
	return dlf.DailyProgram(time.Now())
}

func (dlf Dlf) DailyProgram(day time.Time) ([]Event, error) {
	var url string = dlf.programurl + "?drbm:date=" + day.Format("02.01.2006")
	var resp *http.Response
	var err error
	resp, err = http.Get(url)
	if err != nil {
		return nil, err
	}
	//var rdr *bufio.Reader = bufio.NewReader(resp.Body)

	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(resp.Body)

	var start, title string
	//var current, prev *Event
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		start = s.Find(".time").Text()
		title = s.Find("h3").Text()
		if s.Has("span.title").Length() > 0 {
			fmt.Println("multiple sub entries")
		}
		if strings.LastIndex(title, "aufnehmen") >= 0 {
			title = title[0:strings.LastIndex(title, "aufnehmen")]
		}
		fmt.Printf("%s: %s\n", start, title)
	})

	return nil, err
}
