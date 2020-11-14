package recorder

import "io"
import "net/http"
import "os"
import "time"

type DownloadRecorder struct {
}

func (rec *DownloadRecorder) record(stream string, filepath string, seconds int) error {
	var client http.Client = http.Client{Timeout: time.Duration(seconds) * time.Second}

	var out *os.File
	var err error
	out, err = os.Create(filepath)
	if err != nil {
		return err
	}

	var resp *http.Response
	var err error
	resp, err = client.Get(stream)
	if err != nil {
		out.Close()
		return err
	}

	var n int64
	n, err = io.Copy(resp.Body, out)

	return err
}
