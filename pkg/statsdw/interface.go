package statsdw

type Tag struct {
	Name, Value string
}

type Interface interface {
	Gauge(bucket string, value interface{})
	Increment(bucket string, tags ...Tag)
	Close()
}
