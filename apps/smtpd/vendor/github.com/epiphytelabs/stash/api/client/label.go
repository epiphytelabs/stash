package client

type Label struct {
	Key    string
	Values []string
}

type Labels []Label

func (ls Labels) Get(key string) []string {
	var values []string

	for _, l := range ls {
		if l.Key == key {
			values = append(values, l.Values...)
		}
	}

	return values
}
