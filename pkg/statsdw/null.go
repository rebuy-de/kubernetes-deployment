package statsdw

type NullClient struct{}

func (n NullClient) Gauge(bucket string, value interface{}) {}
func (n NullClient) Increment(bucket string, tags ...Tag)   {}
func (n NullClient) Close()                                 {}
