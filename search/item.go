package search

type item struct {
	number int
	child  *item
	value  *int
}

func (l *item) check(i int) {
	if l.value == nil {
		l.value = &i
		return
	}
	if i > *l.value {
		l.shift(l.value)
		l.value = &i
	} else if l.child != nil {
		l.child.check(i)
	}
}
func (l *item) shift(i *int) {
	if l.child != nil {
		l.child.shift(l.value)
	}
	l.value = i
}
