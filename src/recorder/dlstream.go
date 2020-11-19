package recorder

import "io"
import "net/http"
import "os"
import "time"
import "errors"

type DownloadRecorder struct {
}

const buffer_size uint = 8192

var buf []byte = make([]byte, buffer_size)

func (rec DownloadRecorder) Record(stream string, filepath string, end time.Time) error {
	var client http.Client = http.Client{}

	var out *os.File
	var err error
	out, err = os.Create(filepath)
	if err != nil {
		return err
	}

	var resp *http.Response
	resp, err = client.Get(stream)
	if err != nil {
		out.Close()
		return err
	}

	var read, written int
	for {
		read, err = resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			out.Close()
			resp.Body.Close()
			return err
		}
		if read == 0 {
			// internet connection down for now
			continue
		}
		written, err = out.Write(buf[:read])
		if err != nil {
			out.Close()
			resp.Body.Close()
			return err
		}
		if read != written {
			out.Close()
			resp.Body.Close()
			return errors.New("Could not write file.")
		}

		if time.Now().After(end) {
			return nil
		}
	}

	return err
}
