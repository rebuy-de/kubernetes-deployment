package statsdw

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

var Prefix = "kubernetes-deployment"

type Wrapper struct {
	c *statsd.Client
}

func New(addr string) Interface {
	if addr == "" {
		return NullClient{}
	}

	c, err := statsd.New(
		statsd.Address(addr),
		statsd.Prefix(Prefix),
		statsd.TagsFormat(statsd.Datadog),
	)
	if err != nil {
		log.WithFields(log.Fields{
			"StackTrace": fmt.Sprintf("%+v", errors.WithStack(err)),
		}).Warn("failed to initialize statsd client")
		return NullClient{}
	}

	return &Wrapper{c}
}

func (w *Wrapper) Gauge(bucket string, value interface{}) {
	w.c.Gauge(bucket, value)
}

func (w *Wrapper) Close() {
	w.c.Close()
}

func (w *Wrapper) Increment(bucket string, tags ...Tag) {
	tkv := []string{}
	for _, tag := range tags {
		tkv = append(tkv, tag.Name, tag.Value)
	}

	client := w.c.Clone(
		statsd.Tags(tkv...),
	)

	client.Increment(bucket)
}
