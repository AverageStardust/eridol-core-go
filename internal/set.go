package internal

type Set[V comparable] map[V]struct{}

func (set Set[V]) Add(value V) {
	set[value] = struct{}{}
}

func (set Set[V]) Delete(value V) {
	delete(set, value)
}

func (set Set[V]) Has(value V) (exists bool) {
	_, exists = set[value]
	return
}
