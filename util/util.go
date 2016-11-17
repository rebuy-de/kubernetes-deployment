package util

import (
	"bufio"
	"io"

	log "github.com/Sirupsen/logrus"
)

func PipeToLog(prefix string, rc io.ReadCloser) {
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		log.Debugf("%s %s", prefix, scanner.Text())
	}
}

func PipeToLogrus(ctx log.FieldLogger, rc io.ReadCloser) {
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		ctx.Debug(scanner.Text())
	}
}
