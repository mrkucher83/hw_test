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
		// item := l.items[key].Value.(cacheItem)
		// item.value = value
		l.queue.MoveToFront(l.items[key])
		return true
	}

	l.items[key] = l.queue.PushFront(cacheItem{string(key), value})

	if l.queue.Len() > l.capacity {
		rmKey := l.queue.Back().Value.(cacheItem)
		delete(l.items, Key(rmKey.key))
		l.queue.Remove(l.queue.Back())
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

	for i := l.queue.Back(); i != nil; i = i.Prev {
		l.queue.Remove(i)
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
