package hw04_lru_cache //nolint:golint,stylecheck
import (
	"fmt"
)

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
	key   Key
	Value interface{}
}

func NewCache(capacity int) Cache {
	Cache := lruCache{}
	Cache.capacity = capacity
	Cache.queue = NewList()
	Cache.items = make(map[Key]*ListItem)
	return &Cache
}

func (l lruCache) Set(key Key, value interface{}) bool {
	if _, ok := l.items[key]; ok {
		newCacheItem := cacheItem{}
		newCacheItem.Value = value
		newCacheItem.key = key

		l.items[key].Value = newCacheItem
		i := l.items[key]
		l.queue.MoveToFront(i)
		return true
	}
	newCacheItem := cacheItem{}
	newCacheItem.Value = value
	newCacheItem.key = key

	i := l.queue.PushFront(newCacheItem)
	l.items[key] = i
	if l.queue.Len() > l.capacity {
		var j = l.queue.Back().Value
		i := j.(cacheItem)
		delete(l.items, i.key)
		l.queue.Remove(l.queue.Back())
	}
	return false
}

func (l lruCache) Get(key Key) (interface{}, bool) {
	if _, ok := l.items[key]; ok {
		l.queue.MoveToFront(l.items[key])
		v := l.items[key].Value
		switch i := v.(type) {
		case cacheItem:
			return i.Value, true
		default:
			fmt.Printf("\nUnknown type in func GET %T!\n", v)
		}
	} else {
		return nil, false
	}
	return nil, false
}

func (l lruCache) Clear() {
	for k, v := range l.items {
		var j = l.items[k].Value
		i := j.(cacheItem)
		delete(l.items, i.key)
		l.queue.Remove(v)
	}
}
