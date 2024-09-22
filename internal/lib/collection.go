package lib

import (
	"sync"
)

// Collection is sync.Map, that stores arbitrary data.
type Collection struct {
	items sync.Map
	close chan struct{}
}

// Item just interface for stoing any data.
type item struct {
	data interface{}
}

// NewCollection creates new instance of "collection".
func NewCollection() *Collection {
	c := &Collection{ //nolint:exhaustruct
		close: make(chan struct{}),
	}

	return c
}

// Get returns value by given key. Or empty value.
func (collection *Collection) Get(key interface{}) (interface{}, bool) {
	obj, exists := collection.items.Load(key)

	if !exists {
		return nil, false
	}

	item := obj.(item)

	return item.data, true
}

// Set saves value to given "collection".
func (collection *Collection) Set(key interface{}, value interface{}) {
	collection.items.Store(key, item{
		data: value,
	})
}

// Range apply function f to all keys/values in collection. Retuns immediatly on first error.
func (collection *Collection) Range(f func(key, value interface{}) bool) {
	fn := func(key, value interface{}) bool {
		item := value.(item)

		return f(key, item.data)
	}

	collection.items.Range(fn)
}

// Delete deletes key and corresponding value from collection.
func (collection *Collection) Delete(key interface{}) {
	collection.items.Delete(key)
}

// Close frees memory occupied by collection and deletes collection itself.
func (collection *Collection) Close() {
	collection.close <- struct{}{}
	collection.items = sync.Map{}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
