package stations

import "net/http"
import "time"
import "bufio"
import "fmt"

import "golang.org/x/net/html"

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
	var rdr *bufio.Reader = bufio.NewReader(resp.Body)
	ReadUntilTag(rdr, "span class=\"contenttitle-date\">") // delete up to here
	var program string
	program, err = ReadUntilTag(rdr, "div class=\"dlf-contentright\">")
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	fmt.Println(program)
	return nil, nil
}

func ReadUntilTag(rdr *bufio.Reader, tag string) (string, error) {
	var pre, inner string
	var err error
	var full string = ""
	for {
		pre, err = rdr.ReadString(byte('<'))
		//full += "<"
		if err != nil {
			return full, err
		}
		full += pre
		inner, err = rdr.ReadString(byte('>'))
		//full += ">"
		if err != nil {
			return full, err
		}
		full += inner
		if inner == tag {
			return full, nil
		}
	}
}
