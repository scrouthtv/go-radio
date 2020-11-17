package tui

import "bufio"
import "os"
import "fmt"

var log_writer *bufio.Writer

func init_log() error {
	var log_file *os.File
	var err error
	log_file, err = os.OpenFile("/tmp/radio.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	log_writer = bufio.NewWriter(log_file)
	log_writer.WriteString("Launched\n")
	log_writer.Flush()
	return err
}

func log(msg ...interface{}) {
	if log_writer == nil {
		init_log()
	}
	var word interface{}
	for _, word = range msg {
		fmt.Fprint(log_writer, word)
	}
	log_writer.WriteString("\n")
	log_writer.Flush()
}
