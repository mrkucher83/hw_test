package hw04lrucache

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
	CountItems int // счетчик элементов списка
	FrontItem  *ListItem
	BackItem   *ListItem
}

// реализуем методы интерфейса List для структуры list.
func (l *list) Len() int {
	return l.CountItems
}

func (l *list) Front() *ListItem {
	return l.FrontItem
}

func (l *list) Back() *ListItem {
	return l.BackItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	front := &ListItem{v, l.FrontItem, nil}

	if l.CountItems == 0 {
		l.BackItem = front
	} else {
		front.Next.Prev = front
	}

	l.FrontItem = front
	l.CountItems++

	return front
}

func (l *list) PushBack(v interface{}) *ListItem {
	back := &ListItem{v, nil, l.BackItem}

	if l.CountItems == 0 {
		l.FrontItem = back
	} else {
		back.Prev.Next = back
	}

	l.BackItem = back
	l.CountItems++

	return back
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.FrontItem = i.Next
	} else {
		i.Prev.Next = i.Next
	}
	if i.Next == nil {
		l.BackItem = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.CountItems--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.FrontItem == i {
		return
	}
	if l.BackItem == i {
		l.BackItem = i.Prev
		i.Prev.Next = nil
	} else {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	l.FrontItem.Prev = i
	i.Prev = nil
	i.Next = l.FrontItem
	l.FrontItem = i
}

func NewList() List {
	return new(list)
}
