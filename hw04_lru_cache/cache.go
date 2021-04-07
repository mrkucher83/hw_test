package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
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
		l.items[key].Value = cacheItem{string(key), value}
		l.queue.MoveToFront(l.items[key])
		return true
	}

	l.items[key] = l.queue.PushFront(cacheItem{string(key), value})

	if l.queue.Len() > l.capacity {
		if rmKey, ok := l.queue.Back().Value.(cacheItem); ok {
			delete(l.items, Key(rmKey.key))
			l.queue.Remove(l.queue.Back())
		}
	}
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	if _, ok := l.items[key]; ok {
		l.queue.MoveToFront(l.items[key])
		val := l.items[key].Value.(cacheItem)
		return val.value, true
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
