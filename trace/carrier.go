package trace

type Carrier interface {
	Get(key string) string
	Set(key, value string)
}

type MapCarrier map[string]string

func (c MapCarrier) Get(key string) string {
	return c[key]
}

func (c MapCarrier) Set(key, value string) {
	c[key] = value
}
