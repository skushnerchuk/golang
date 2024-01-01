package hw04lrucache

import "sync"

type Key string

type queueItems map[Key]*ListItem

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    queueItems
	lock     sync.Mutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(queueItems, capacity),
	}
}

// Set Обновление элемента или установка нового в кэш.
func (l *lruCache) Set(key Key, v interface{}) bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	// Если такой элемент уже есть в кэше, перемещаем его на первое место
	if _, ok := l.items[key]; ok {
		l.queue.MoveToFront(l.items[key])
		l.queue.Front().Value = &cacheItem{key: key, value: v}
		l.items[key].Value = l.queue.Front().Value
		return true
	}

	// Если пытаемся добавить новый элемент с превышением ёмкости
	// кэша - отстреливаем из него самый старый элемент
	if l.queue.Len() == l.capacity {
		delete(l.items, l.queue.Back().Value.(*cacheItem).key)
		l.queue.Remove(l.queue.Back())
	}

	// Помещаем новый элемент в голову кэша
	l.items[key] = l.queue.PushFront(&cacheItem{key: key, value: v})
	return false
}

// Get Получение элемента из кэша по ключу.
func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.lock.Lock()
	defer l.lock.Unlock()

	// Если элемент есть в кэше - перемещаем его в голову и возвращаем значение
	if _, ok := l.items[key]; ok {
		l.queue.MoveToFront(l.items[key])
		l.items[key] = l.queue.Front()
		res := l.items[key].Value.(*cacheItem).value
		return res, true
	}
	return nil, false
}

// Clear Очистка кэша.
func (l *lruCache) Clear() {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.queue = NewList()
	l.items = make(queueItems, l.capacity)
}
