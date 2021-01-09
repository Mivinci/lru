package lru

import "container/list"

type Cache struct {
	cap   int
	Evict func(k interface{}, v interface{})
	l     *list.List
	cache map[interface{}]*list.Element
}

type entry struct {
	key   interface{}
	value interface{}
}

func New(cap int) *Cache {
	return &Cache{
		cap:   cap,
		l:     list.New(),
		cache: make(map[interface{}]*list.Element),
	}
}

func (c *Cache) Add(key, value interface{}) {
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.l = list.New()
	}
	if el, ok := c.cache[key]; ok {
		c.l.MoveToFront(el)
		el.Value.(*entry).value = value
		return
	}
	el := c.l.PushFront(&entry{key, value})
	c.cache[key] = el
	if c.cap != 0 && c.l.Len() > c.cap {
		c.RemoveOldest()
	}
}

func (c *Cache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	el := c.l.Back()
	if el != nil {
		c.removeElement(el)
	}
}

func (c *Cache) removeElement(e *list.Element) {
	c.l.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
	if c.Evict != nil {
		c.Evict(kv.key, kv.value)
	}
}

func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	if el, ok := c.cache[key]; ok {
		c.l.MoveToFront(el)
		return el.Value.(*entry).value, true
	}
	return
}

func (c *Cache) Remove(key interface{}) {
	if c.cache == nil {
		return
	}
	if el, ok := c.cache[key]; ok {
		c.removeElement(el)
	}
}

func (c *Cache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.l.Len()
}

func (c *Cache) Clear() {
	if c.Evict != nil {
		for _, e := range c.cache {
			kv := e.Value.(*entry)
			c.Evict(kv.key, kv.value)
		}
	}
	c.l = nil
	c.cache = nil
}

func (c *Cache) Walk(fn func(k interface{}, v interface{}) error) error {
	for k, v := range c.cache {
		if err := fn(k, v.Value.(*entry).value); err != nil {
			return err
		}
	}
	return nil
}
