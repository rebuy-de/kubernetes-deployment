package util

import (
	"bufio"
	"io"
	"log"
)

func PipeToLog(prefix string, rc io.ReadCloser) {
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		log.Printf("%s %s", prefix, scanner.Text())
	}
}