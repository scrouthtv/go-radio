package stations

import "net/http"
import "time"
import "bufio"
import "fmt"

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
	ReadUntilTag(rdr, []rune("span class=\"contenttitle-date\">")) // delete up to here
	panic("last")
	var program string
	program, err = ReadUntilTag(rdr, []rune("div class=\"dlf-contentright\">"))
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	fmt.Println(program)
	return nil, nil
}

func ReadUntilTag(rdr *bufio.Reader, tag []rune) (string, error) {
	var read string
	var fullRead string = ""
	var err error
	var i int
	var taglen int = len(tag)
	var readRune rune
	var validTag bool = true
	var thisTag string
	for {
		read, err = rdr.ReadString(byte('<'))
		fullRead = fullRead + read
		validTag = true
		thisTag = ""
		for i = 0; i < taglen; i++ {
			readRune, _, err = rdr.ReadRune()
			if err != nil {
				fmt.Println("reading <")
				return "", err
			}
			read += string(readRune)
			thisTag += string(readRune)
			if readRune == '>' {
				validTag = false
				break
			} else if validTag && readRune != tag[i] {
				validTag = false
				break
			}
		}
		if validTag {
			return read, nil
		} else {
			fullRead = fullRead + thisTag
			if i > 0 {
				//fmt.Println("discarding", thisTag)
				//fmt.Println(fullRead)
				fmt.Println(read)
				fmt.Println("--------------")
			}
		}
	}
}
