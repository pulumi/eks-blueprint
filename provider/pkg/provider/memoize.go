package provider

import (
	"sync"
	"time"

	"github.com/kofalt/go-memoize"
)

var lock = &sync.Mutex{}

var memoCache *memoize.Memoizer

func getMemoCache() *memoize.Memoizer {
	if memoCache == nil {
		lock.Lock()
		defer lock.Unlock()
		if memoCache == nil {
			memoCache = memoize.NewMemoizer(90*time.Second, 10*time.Minute)
		}
	}
	return memoCache
}
