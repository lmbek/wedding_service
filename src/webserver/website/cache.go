package website

import (
	"fmt"
	"sync"
)

var cacheMu sync.RWMutex
var pageCache = make(map[string]string)

func ReadCache(templateCachingEnabled bool, path string) (string, bool) {
	if !templateCachingEnabled {
		return "", false
	}
	fmt.Println("read cache")
	cacheMu.RLock()
	page, ok := pageCache[path]
	cacheMu.RUnlock()
	return page, ok
}

func UpdateCache(templateCachingEnabled bool, path string, page string) bool {
	if !templateCachingEnabled {
		return false
	}
	fmt.Println("updated cache")
	cacheMu.Lock()
	pageCache[path] = page
	cacheMu.Unlock()
	return true
}
