package cmd

import (
	log "github.com/Sirupsen/logrus"
)

const ()

var ()

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{})
}
