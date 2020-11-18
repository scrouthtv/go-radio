package stations

import "net/http"
import "time"
import "strings"
import "net/url"

import "github.com/PuerkitoBio/goquery"

type Dlf struct {
	name         string
	url          string
	programurl   string
	programStart string // some tag before or at the start of the program table
	subCatIdent  string // some tag that only occurs inside events that have multiple sub events
	subTimeIdx   int
	subTitleIdx  int
	linksInTitle int
	linksInSub   int // how many <a>s are in the title in a sub, 1 for Dlf since they have the record button
}

// scrap the note on linksInSub. both pages have the record button only for some events
var Deutschlandfunk Dlf = Dlf{"Deutschlandfunk", "https://st01.sslstream.dlf.de/dlf/01/high/aac/stream.aac", "https://www.deutschlandfunk.de/programmvorschau.281.de.html", ".dlf-contentleft", "td.description", 0, 1, 1, 0}
var DlfKultur Dlf = Dlf{"Deutschlandfunk", "https://st02.sslstream.dlf.de/dlf/02/high/aac/stream.aac", "https://www.deutschlandfunkkultur.de/programmvorschau.282.de.html", ".drk-tagesprogramm", "p.subDescription", 1, 2, 1, 0}

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
	var htmlurl string = dlf.programurl + "?drbm:date=" + day.Format("02.01.2006")
	var resp *http.Response
	var err error
	resp, err = http.Get(htmlurl)
	if err != nil {
		return nil, err
	}

	var events []Event

	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(resp.Body)

	// TODO return slice of all errors that occured

	var start, title, desc, category, url, caturl string
	var endTime time.Time
	var current *Event
	doc.Find(dlf.programStart).Find("tr").Each(func(i int, s *goquery.Selection) {
		if s.Find(dlf.subCatIdent).Has("span.title").Length() > 0 { // has multiple sub events
			start = s.Find("td").Slice(0, 1).Text()
			if len(events) > 0 {
				endTime, err = time.Parse("15:04 Uhr", start)
				endTime = combine(day, endTime, combineMaskYMD)
				events[len(events)-1].End = endTime
			}
			category = s.Find("h3").Slice(0, 1).Text()
			if s.Find("h3").Find("a").Length() > dlf.linksInTitle {
				caturl, _ = s.Find("h3").Find("a").Attr("href")
			}
			current = createEvent(&dlf, day, &category, &desc, &start, nil, &category, &url, &caturl)
			events = append(events, *current)
			s.Find("p").Each(func(i int, sub *goquery.Selection) {
				if sub.HasClass("subDescription") {
					sub.Find("span").Each(func(i int, span *goquery.Selection) {
						if i == dlf.subTimeIdx { // the first element is the start time
							start = span.Text()
							if len(events) > 0 {
								endTime, err = time.Parse("15:04", start[0:5])
								endTime = combine(day, endTime, combineMaskYMD)
								events[len(events)-1].End = endTime
							}
						} else if i == dlf.subTitleIdx { // the second is the title
							title = span.Text()
							if span.Find("a").Length() > dlf.linksInSub {
								url, _ = span.Find("a").Slice(0, 1).Attr("href")
							}

							current = createEvent(&dlf, day, &title, &desc, &start, nil, &category, &url, &caturl)
							events = append(events, *current)
							// these are all the sub events in a multievent
						}
					})
				} else {
					var descHtml string
					descHtml, err = sub.Html()
					events[len(events)-1].Info += formatMultilineHtml(descHtml)
				}
			})
		} else {
			category = "N/A"
			caturl = ""
			s.Find("td").Each(func(i int, cell *goquery.Selection) {
				if cell.HasClass("time") {
					start = cell.Text()
					if len(events) > 0 {
						endTime, err = time.Parse("15:04", start[0:5])
						endTime = combine(day, endTime, combineMaskYMD)
						events[len(events)-1].End = endTime
					}
				} else {
					title = s.Find("h3").Text()
					if s.Find("a").Length() > dlf.linksInTitle {
						url, _ = s.Find("h3").Find("a").Slice(0, 1).Attr("href")
						if strings.HasPrefix(url, "deutschlandradio-recorder-programmieren.1406.de.html") {
							url = ""
						}
					}
					desc = ""
					var paraHtml string
					cell.Find("p").Each(func(i int, paragraph *goquery.Selection) {
						paraHtml, err = paragraph.Html()
						desc += paraHtml + "\n"
					})
					desc = formatMultilineHtml(desc)
				}
			})
			current = createEvent(&dlf, day, &title, &desc, &start, nil, &category, &url, &caturl)
			events = append(events, *current)
		}
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

func createEvent(dlf *Dlf, day time.Time, name *string, info *string, start *string, end *string, category *string, url *string, caturl *string) *Event {
	var startTime, endTime time.Time
	if start != nil {
		startTime, _ = time.Parse("15:04", trimSpaces(*start)[0:5])
		startTime = combine(day, startTime, combineMaskYMD)
	}
	if end != nil {
		endTime, _ = time.Parse("15:04", trimSpaces(*end))
		endTime = combine(day, endTime, combineMaskYMD)
		panic("not tested")
	}

	var ev Event = Event{trimSpaces(*name), strings.Trim(*info, "\n "), startTime, endTime, trimSpaces(*category),
		absUrl(dlf, *url), absUrl(dlf, *caturl)}
	var emptyString string = ""
	name = &emptyString
	info = &emptyString
	start = &emptyString
	end = &emptyString
	category = &emptyString
	url = &emptyString
	caturl = &emptyString

	return &ev
}

func absUrl(dlf *Dlf, urlStr string) url.URL {
	var myurl *url.URL
	myurl, _ = url.Parse(urlStr)

	if urlStr == "" {
		return url.URL{}
	}

	if myurl.IsAbs() {
		return *myurl
	} else {
		var base *url.URL
		base, _ = url.Parse(dlf.programurl)
		myurl = base.ResolveReference(myurl)
		return *myurl
	}
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
