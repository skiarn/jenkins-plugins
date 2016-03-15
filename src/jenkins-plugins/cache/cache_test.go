package cache

import (
	"bytes"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

//go test -bench=.
//go test cache/ -run=10 -bench=.
var (
	data = []byte("This is data to be stored...")
	key  = "key"
)

func initCache(size int) (c *Cache) {
	cache := New()
	// run the set and get function b.N times
	for i := 0; i < size; i++ {
		var keyBuffer bytes.Buffer
		keyBuffer.WriteString(key)
		keyBuffer.WriteString(strconv.Itoa(size))
		cache.Set(keyBuffer.String(), data)
	}
	return &cache
}
func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func benchmarkCache(cacheSizeToTest int, b *testing.B) {
	cache := initCache(cacheSizeToTest)
	tempKey := key + strconv.Itoa(random(0, cacheSizeToTest))
	// run the set and get function b.N times
	for n := 0; n < b.N; n++ {
		done := make(chan bool)
		for x := 0; x < cacheSizeToTest; x++ {
			go func() {
				cache.Set(tempKey, data)
				_ = cache.Get(tempKey)
				done <- true
			}()
		}

	}

}

func BenchmarkCache100(b *testing.B)      { benchmarkCache(100, b) }
func BenchmarkCache100000(b *testing.B)   { benchmarkCache(100000, b) }
func BenchmarkCache1000000(b *testing.B)  { benchmarkCache(1000000, b) }
func BenchmarkCache10000000(b *testing.B) { benchmarkCache(2000000, b) }

func TestCacheSize(t *testing.T) {
	expectedSize := 1000
	cache := New()
	for n := 0; n < expectedSize; n++ {
		var keyBuffer bytes.Buffer
		keyBuffer.WriteString(key)
		keyBuffer.WriteString(strconv.Itoa(n))
		cache.Set(keyBuffer.String(), data)
		_ = cache.Get(keyBuffer.String())
	}

	if result := cache.Size(); result != expectedSize {
		t.Errorf("Size() returned %t, expected: %t", result, expectedSize)
	}
}

func TestCacheGetAll(t *testing.T) {
	cache := New()
	cache.Set("keyone", data)
	cache.Set("keytwo", data)

	keys := cache.GetAll()

	existing := make(map[string]bool)
	for _, v := range keys {
		existing[v] = true
	}

	if !stringInSlice("keyone", keys) || !stringInSlice("keytwo", keys) {
		t.Errorf("Expected keys %s in cache.", keys)
	}
}
func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
