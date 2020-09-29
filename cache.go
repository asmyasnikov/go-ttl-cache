package cache

import (
	"regexp"
	"sync"
	"time"
)

// Item is a cached reference
type Item struct {
	Content    interface{}
	Expiration time.Time
}

// Expired returns true if the item has expired.
func (item Item) Expired() bool {
	if item.Expiration.IsZero() {
		return false
	}
	return time.Now().After(item.Expiration)
}

//Storage mecanism for caching strings in memory
type Storage struct {
	items map[string]Item
	mtx   *sync.RWMutex
}

//NewStorage creates a new in memory storage
func NewStorage() *Storage {
	return &Storage{
		items: make(map[string]Item),
		mtx:   &sync.RWMutex{},
	}
}

//Get a cached content by key
func (s Storage) Get(key string) interface{} {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	item, ok := s.items[key]
	if !ok {
		return nil
	}
	if item.Expired() {
		delete(s.items, key)
		return nil
	}
	return item.Content
}

//Expiration of item by key
func (s Storage) Expiration(key string) time.Time {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	item, ok := s.items[key]
	if !ok {
		return time.Time{}
	}
	return item.Expiration
}

//TTL of item by key
func (s Storage) TTL(key string) time.Duration {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	item, ok := s.items[key]
	if !ok {
		return 0
	}
	return time.Until(item.Expiration)
}

//Get a cached content by key
func (s Storage) Keys(re string) []string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	keys := make([]string, 0)
	matcher := regexp.MustCompile(re)

	for k := range s.items {
		if matcher.MatchString(k) {
			keys = append(keys, k)
		}
	}

	return keys
}

//Set a cached content by key
func (s Storage) Set(key string, content interface{}, duration time.Duration) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.items[key] = Item{
		Content: content,
		Expiration: func() time.Time {
			if duration == 0 {
				return time.Time{}
			}
			return time.Now().Add(duration)
		}(),
	}
}

//Flush make clear cache starts with prefix
func (s Storage) Flush(prefix string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	for k := range s.items {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			delete(s.items, k)
		}
	}
}

//Rem make clear cache starts with prefix
func (s Storage) Rem(key string) interface{} {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if v, ok := s.items[key]; ok {
		delete(s.items, key)
		return v.Content
	}
	return nil
}
