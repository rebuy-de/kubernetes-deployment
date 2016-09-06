package util

import (
	"bufio"
	"bytes"
	"io"
	"log"
)

func PipeToLog(prefix string, rc io.ReadCloser) {
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		log.Printf("%s %s", prefix, scanner.Text())
	}
}

func PipeToString(rc io.ReadCloser, c chan string) string{
	scanner := bufio.NewScanner(rc)
	var buffer bytes.Buffer
	for scanner.Scan() {
		buffer.WriteString(scanner.Text())
	}
	c <- buffer.String()
	return buffer.String()
}
