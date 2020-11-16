package stations

import "net/http"
import "time"
import "bufio"
import "fmt"
import "strings"

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

	//var events []Event
	var row string
	rdr = bufio.NewReader(strings.NewReader(program))
	var rowRdr *bufio.Reader

	//var name, info, start, duration, category string
	var start string
	for {
		fmt.Println("loop")
		ReadUntilTag(rdr, "tr id=\"anc")
		row, err = ReadUntilTag(rdr, "/tr")
		if err != nil {
			break // should be continue
		}

		// for now I will assume a fixed order of the fields in the table
		// and do this performance-wise in the worst possible way:
		rowRdr = bufio.NewReader(strings.NewReader(row))
		_, err = ReadUntilTag(rowRdr, "td class=\"time\">")
		if err != nil {
			fmt.Println("err 1", err)
			break
		}
		start, err = ReadUntilTag(rowRdr, "") // up to the closing bracket
		if err != nil {
			fmt.Println("err 2", err)
			break
		}
		fmt.Println(start)
		var i int
		for i = 0; i < 20; i++ {
			fmt.Println()
		}
	}

	return nil, nil
}

func ReadUntilTag(rdr *bufio.Reader, tag string) (string, error) {
	var pre, inner string
	var err error
	var full strings.Builder
	for {
		pre, err = rdr.ReadString(byte('<'))
		if err != nil {
			return full.String(), err
		}
		full.WriteString(pre)
		fmt.Println(inner)
		fmt.Println(tag)
		fmt.Println("---")
		inner, err = rdr.ReadString(byte('>'))
		if err != nil {
			return full.String(), err
		}
		full.WriteString(inner)
		if strings.HasPrefix(inner, tag) {
			return full.String(), nil
		}
	}
}
