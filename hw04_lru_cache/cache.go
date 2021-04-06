package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	Cache // Remove me after realization.

	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   string
	value interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	if _, ok := l.items[key]; ok {
		l.items[key].Value = value
		l.queue.MoveToFront(l.items[key])
		return true
	}

	l.items[key] = l.queue.PushFront(value)
	if l.queue.Len() > l.capacity {
		l.queue.Remove(l.queue.Back())
		// todo: удалить значение из словаря
	}
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	if _, ok := l.items[key]; ok {
		l.queue.MoveToFront(l.items[key])
		return l.items[key].Value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.items = nil
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
