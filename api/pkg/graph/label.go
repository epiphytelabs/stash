package graph

type Label struct {
	key    string
	values []string
}

func (l *Label) Key() string {
	return l.key
}

func (l *Label) Values() []string {
	return l.values
}
