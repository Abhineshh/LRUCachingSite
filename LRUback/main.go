package main

import (
	"fmt"
	"net/http"
	//"strconv"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type entry struct {
	value      interface{}
	expiration time.Time
}

type LRUCache struct {
	cache    map[string]*entry
	lruList  []string
	mutex    sync.Mutex
	capacity int
}

func LRUinitializer(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*entry),
		lruList:  make([]string, 0, capacity),
	}
}

func (c *LRUCache) Set(key string, value interface{}, expiration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.cache) >= c.capacity {
		delete(c.cache, c.lruList[len(c.lruList)-1])
		c.lruList = c.lruList[:len(c.lruList)-1]
	}

	c.cache[key] = &entry{
		value:      value,
		expiration: time.Now().Add(expiration),
	}
	c.lruList = append([]string{key}, c.lruList...)
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if entry, ok := c.cache[key]; ok {
		if entry.expiration.After(time.Now()) {
			// Move key to front of LRU list
			for i, k := range c.lruList {
				if k == key {
					c.lruList = append(append([]string{key}, c.lruList[:i]...), c.lruList[i+1:]...)
					break
				}
			}
			return entry.value, true
		}
		// Expired entry, delete from cache
		delete(c.cache, key)
		for i, k := range c.lruList {
			if k == key {
				c.lruList = append(c.lruList[:i], c.lruList[i+1:]...)
				break
			}
		}
	}
	return nil, false
}

func main() {
	cache := LRUinitializer(1024)

	router := mux.NewRouter()

	router.HandleFunc("/get/{key}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]
		if value, ok := cache.Get(key); ok {
			fmt.Fprintf(w, "%v", value)
		} else {
			fmt.Println("ding dong")
			http.NotFound(w, r)
		}
	}).Methods("GET")

	router.HandleFunc("/set/{key}/{value}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]
		value := vars["value"]
		cache.Set(key, value, 5*time.Second)
		fmt.Fprintf(w, "Key %s set with value %s", key, value)
	}).Methods("POST")

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}), // Add your React app's origin here
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	http.ListenAndServe(":8585", cors(router))
}
