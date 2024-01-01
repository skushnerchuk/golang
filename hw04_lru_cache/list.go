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

// Len Получение длины списка
func (l *list) Len() int {
	return l.length
}

// Front Получение ссылки на головной элемент списка
func (l *list) Front() *ListItem {
	return l.head
}

// Back Получение ссылки на последний элемент списка
func (l *list) Back() *ListItem {
	return l.tail
}

// PushFront Размещение элемента в начале списка
func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.head != nil {
		item.Next = l.head
		l.head.Prev = item
	} else {
		l.tail = item
	}
	l.head = item
	l.length++
	return item
}

// PushBack Размещение элемента в конце списка
func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.tail != nil {
		item.Prev = l.tail
		l.tail.Next = item
	} else {
		l.head = item
	}
	l.tail = item
	l.length++
	return item
}

// Remove Удаление элемента из списка
func (l *list) Remove(item *ListItem) {
	if item == nil || l.length == 0 {
		return
	}

	// Если удаляется первый элемент списка - говорим что следующий
	// элемент становится первым, иначе замыкаем соседние элементы друг на друга
	if item.Prev == nil {
		l.head = item.Next
		item.Next.Prev = nil
	} else {
		item.Prev.Next = item.Next
	}

	// Если удаляется последний элемент списка - говорим что предыдущий
	// элемент становится последним, иначе замыкаем соседние элементы друг на друга
	if item.Next == nil {
		l.tail = item.Prev
		item.Prev.Next = nil
	} else {
		item.Next.Prev = item.Prev
	}

	l.length--
}

// MoveToFront Перемещение элемента в начало списка
func (l *list) MoveToFront(item *ListItem) {
	// Если передан существующий объект и это не голова списка, то выполняем его перемещение
	if item != nil && l.head != item {
		itemValue := item.Value
		l.Remove(item)
		l.PushFront(itemValue)
	}
}

// NewList Создание нового списка
func NewList() List {
	return new(list)
}
