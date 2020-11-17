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

	var events []Event

	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(resp.Body)

	if true == false {
		fmt.Println("dont annoy me w unused imports if im trying to debug ffs")
		strings.Split("asdf", "")
	}

	var start, title, desc, category string
	var current *Event
	// dlf-contentleft might be named differently on other sites
	doc.Find(".dlf-contentleft").Find("tr").Each(func(i int, s *goquery.Selection) {
		if s.Find("td.description").Has("span.title").Length() > 0 { // has multiple sub events
			start = s.Find("td").Slice(0, 1).Text()
			if len(events) > 0 {
				events[len(events)-1].End, err = time.Parse("15:04 Uhr", start)
			}
			category = s.Find("h3").Slice(0, 1).Text()
			current = createEvent(&category, &desc, &start, nil, &category)
			events = append(events, *current)
			s.Find("p").Each(func(i int, sub *goquery.Selection) {
				if sub.HasClass("subDescription") {
					if start == "10:08 Uhr" {
						/*fmt.Println("no info:")
						fmt.Println(sub.Text())*/
					}
					sub.Find("span").Each(func(i int, span *goquery.Selection) {
						if i == 0 { // the first element is the start time
							start = span.Text()
							if len(events) > 0 {
								events[len(events)-1].End, err = time.Parse("15:04", start)
							}
						} else { // the second is the title
							title = span.Text()
							current = createEvent(&title, &desc, &start, nil, &category)
							events = append(events, *current)
							// these are all the sub events in a multievent
						}
					})
				} else {
					if start == "10:08 Uhr" {
						/*fmt.Println("info")
						fmt.Println(sub.Text())*/
					}
					var descHtml string
					descHtml, err = sub.Html()
					events[len(events)-1].Info += formatMultilineHtml(descHtml)
				}
			})
		} else {
			category = "N/A"
			s.Find("td").Each(func(i int, cell *goquery.Selection) {
				if cell.HasClass("time") {
					start = cell.Text()
					if len(events) > 0 {
						events[len(events)-1].End, err = time.Parse("15:04 Uhr", start)
					}
				} else {
					title = s.Find("h3").Text()
					desc = ""
					var paraHtml string
					cell.Find("p").Each(func(i int, paragraph *goquery.Selection) {
						paraHtml, err = paragraph.Html()
						desc += paraHtml + "\n"
					})
					desc = formatMultilineHtml(desc)
				}
			})
			current = createEvent(&title, &desc, &start, nil, &category)
			events = append(events, *current)
		}
		//fmt.Printf("%s: %s\n", start, title)
	})

	return events, err
}

func formatMultilineHtml(html string) string {
	html = strings.ReplaceAll(html, "<br>", "\n")
	html = strings.ReplaceAll(html, "<br/>", "\n")
	html = strings.ReplaceAll(html, "<br />", "\n")
	html = strings.ReplaceAll(html, "</p>", "\n")
	var doc *goquery.Document
	var err error
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err == nil {
		return doc.Text()
	} else {
		return html
	}
}

func createEvent(name *string, info *string, start *string, end *string, category *string) *Event {
	var startTime, endTime time.Time
	if start != nil {
		if strings.HasSuffix(*start, " Uhr") {
			startTime, _ = time.Parse("15:04 Uhr", *start)
		} else {
			startTime, _ = time.Parse("15:04", *start)
		}
	}
	if end != nil {
		endTime, _ = time.Parse("15:04", *end)
		panic("not tested")
	}

	var ev Event = Event{trimSpaces(*name), strings.Trim(*info, "\n "), startTime, endTime, trimSpaces(*category)}
	var emptyString string = ""
	name = &emptyString
	info = &emptyString
	start = &emptyString
	end = &emptyString
	category = &emptyString

	return &ev
}

func trimSpaces(str string) string {
	str = strings.ReplaceAll(str, "&nbsp;", "")
	str = strings.Trim(str, " \t	  ")
	if strings.HasSuffix(str, "aufnehmen") {
		str = str[0:strings.LastIndex(str, "aufnehmen")]
	}
	str = strings.Trim(str, " \t	  ")
	return str
}
