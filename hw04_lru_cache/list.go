package hw04_lru_cache //nolint:golint,stylecheck

// var ErrEmptyList = errors.New("trying to remove an item from empty list")

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newListItem := &ListItem{
		Value: v,
		Next:  l.front,
		Prev:  nil,
	}
	if l.len > 0 {
		l.front.Prev = newListItem
	} else {
		l.back = newListItem
	}
	l.front = newListItem
	l.len++
	return newListItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newListItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.back,
	}
	if l.len > 0 {
		l.back.Next = newListItem
	} else {
		l.front = newListItem
	}
	l.back = newListItem
	l.len++
	return newListItem
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.PushFront(i.Value)
	l.Remove(i)
}

func NewList() List {
	return &list{}
}
