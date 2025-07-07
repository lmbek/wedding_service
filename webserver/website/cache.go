package website

import (
	"fmt"
	"sync"
	"wedding_service/flags"
)

var cacheMu sync.RWMutex
var pageCache = make(map[string]string)

func ReadCache(path string) (string, bool) {
	if flags.FrontendFlag != "" {
		return "", false
	}
	fmt.Println("read cache")
	cacheMu.RLock()
	page, ok := pageCache[path]
	cacheMu.RUnlock()
	return page, ok
}

func UpdateCache(path string, page string) bool {
	if flags.FrontendFlag != "" {
		return false
	}
	fmt.Println("updated cache")
	cacheMu.Lock()
	pageCache[path] = page
	cacheMu.Unlock()
	return true
}
