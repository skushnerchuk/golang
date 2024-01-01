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
	length int
	head   *ListItem
	tail   *ListItem
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	newValue := new(ListItem)
	newValue.Value = v
	if l.tail == nil && l.head == nil {
		l.head = newValue
		l.tail = newValue
	} else {
		newValue.Next = l.head
		l.head = newValue
	}
	l.length++
	return newValue
}

func (l *list) PushBack(v interface{}) *ListItem {
	newValue := new(ListItem)
	newValue.Value = v
	if l.tail == nil && l.head == nil {
		l.head = newValue
		l.tail = newValue
	} else {
		newValue.Prev = l.tail
		l.tail.Next = newValue
		l.tail = newValue
	}
	l.length++
	return newValue
}

func (l *list) Remove(item *ListItem) {
	if item == nil || l.length == 0 {
		return
	}
	prev := item.Prev
	next := item.Next
	switch {
	// Если элемент всего один в списке - обнуляем указатели на начало и конец
	case prev == nil && next == nil:
		l.tail = nil
		l.head = nil
	// Если это первый в списке элемент, переключаем указатель начала на следующий
	// и обнуляем указатель на удаляемый элемент
	case prev == nil:
		l.head = next
		next.Prev = nil
	// Если это последний элемент в списке, переключаем указатель на предыдущий и обнуляем в нем
	// указатель на удаляемый элемент
	case next == nil:
		l.tail = prev
		prev.Next = nil
	// Если удаляем из середины списка, то замыкаем указатели соседних элементов друг на друга
	default:
		prev.Next = next
		next.Prev = prev
	}
	l.length--
}

func (l *list) MoveToFront(item *ListItem) {
	if item != nil {
		itemValue := item.Value
		l.Remove(item)
		l.PushFront(itemValue)
	}
}

func NewList() List {
	return new(list)
}
