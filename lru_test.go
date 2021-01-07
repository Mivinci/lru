package lru

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func True(t *testing.T, ok bool) {
	assert(t, ok, 1)
}

func Equal(t *testing.T, expect, got interface{}) {
	assert(t, reflect.DeepEqual(expect, got), 1)
}

func assert(t *testing.T, ok bool, cd int) {
	if !ok {
		_, file, line, _ := runtime.Caller(cd + 1)
		t.Errorf("%s:%d", file, line)
		t.FailNow()
	}
}

func TestGet(t *testing.T) {
	c := New(3)
	c.Add("a", 1)
	Equal(t, 1, c.Len())
	v, ok := c.Get("a")
	True(t, ok)
	Equal(t, 1, v)
	c.Add("b", 2)
	Equal(t, 2, c.Len())
	c.Add("c", 3)
	c.Add("d", 4)
	vv, ok := c.Get("a")
	True(t, !ok)
	Equal(t, nil, vv)
	c.Remove("d")
	Equal(t, 2, c.Len())
	Equal(t, 2, c.l.Back().Value.(*entry).value)
}

func Example() {
	c := New(3)
	c.Evict = func(k, v interface{}) {
		fmt.Println(k, v)
	}
	c.Add("a", 1)
	c.Remove("a")

	// Output:
	// a 1
}
