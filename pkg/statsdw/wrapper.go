package statsdw

import statsd "gopkg.in/alexcesaro/statsd.v2"

var Prefix = "kubernetes-deployment"

type Wrapper struct {
	c *statsd.Client
}

func New(addr string) (Interface, error) {
	if addr == "" {
		return NullClient{}, nil
	}

	c, err := statsd.New(
		statsd.Address(addr),
		statsd.Prefix(Prefix),
		statsd.TagsFormat(statsd.Datadog),
	)
	if err != nil {
		return nil, err
	}

	return &Wrapper{c}, nil
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
