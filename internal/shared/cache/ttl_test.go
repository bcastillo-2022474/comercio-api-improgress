package cache

import (
	"sync"
	"testing"
	"time"
)

type fakeClock struct {
	mu sync.Mutex
	t  time.Time
}

func (f *fakeClock) Now() time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.t
}

func (f *fakeClock) Advance(d time.Duration) {
	f.mu.Lock()
	f.t = f.t.Add(d)
	f.mu.Unlock()
}

func TestTTLGetSetAndExpiry(t *testing.T) {
	clk := &fakeClock{t: time.Unix(1000, 0)}
	c := NewTTL[string, int](time.Minute, clk.Now)

	if _, ok := c.Get("a"); ok {
		t.Fatal("empty cache returned a hit")
	}
	c.Set("a", 1)
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("got %v,%v want 1,true", v, ok)
	}
	clk.Advance(59 * time.Second)
	if _, ok := c.Get("a"); !ok {
		t.Fatal("entry expired early")
	}
	clk.Advance(time.Second)
	if _, ok := c.Get("a"); ok {
		t.Fatal("entry should be expired at ttl")
	}
}

func TestTTLOverwriteRefreshesExpiry(t *testing.T) {
	clk := &fakeClock{t: time.Unix(1000, 0)}
	c := NewTTL[string, int](time.Minute, clk.Now)
	c.Set("a", 1)
	clk.Advance(40 * time.Second)
	c.Set("a", 2)
	clk.Advance(40 * time.Second)
	if v, ok := c.Get("a"); !ok || v != 2 {
		t.Fatalf("got %v,%v want 2,true", v, ok)
	}
}

func TestTTLSetSweepsExpired(t *testing.T) {
	clk := &fakeClock{t: time.Unix(1000, 0)}
	c := NewTTL[int, int](time.Minute, clk.Now)
	for i := 0; i < 10; i++ {
		c.Set(i, i)
	}
	clk.Advance(2 * time.Minute)
	c.Set(100, 100)
	if n := len(c.items); n != 1 {
		t.Fatalf("map has %d entries after sweep, want 1", n)
	}
}

func TestTTLConcurrent(t *testing.T) {
	c := NewTTL[int, int](time.Millisecond, nil)
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < 500; i++ {
				c.Set(i%10, g)
				c.Get(i % 10)
			}
		}(g)
	}
	wg.Wait()
}
