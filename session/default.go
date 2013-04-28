package session

import (
	"github.com/sdming/mcache"
	"time"
)

type sessionEntry map[string]interface{}

type MCache struct {
	cache *mcache.MCache
}

func newMCache() *MCache {
	return &MCache{
		cache: mcache.NewMCache(),
	}
}

func (mc *MCache) Name() string {
	return "default"
}

func (mc *MCache) get(id string) (sessionEntry, bool) {
	if id == "" {
		return nil, false
	}
	x, ok := mc.cache.Get(id)
	return x.(sessionEntry), ok
}

func (mc *MCache) Add(sessionId, key string, value interface{}) (bool, error) {
	entry, ok := mc.get(sessionId)
	if !ok || entry == nil {
		return false, SessionNotExists
	}
	if _, ok := entry[key]; ok {
		return false, nil
	}
	entry[key] = value
	return true, nil
}

func (mc *MCache) Get(sessionId, key string) (interface{}, bool, error) {
	entry, ok := mc.get(sessionId)
	if !ok || entry == nil {
		return nil, false, SessionNotExists
	}
	x, ok := entry[key]
	return x, ok, nil
}

func (mc *MCache) Set(sessionId, key string, value interface{}) error {
	entry, ok := mc.get(sessionId)
	if !ok || entry == nil {
		return SessionNotExists
	}
	entry[key] = value
	return nil
}
func (mc *MCache) Remove(sessionId, key string) error {
	entry, ok := mc.get(sessionId)
	if !ok || entry == nil {
		return SessionNotExists
	}
	delete(entry, key)
	return nil
}

func (mc *MCache) New(sessionId string, timeout time.Duration) error {
	entry := make(sessionEntry)
	mc.cache.SetSlid(sessionId, entry, timeout)
	return nil
}

func (mc *MCache) Abandon(sessionId string) error {
	mc.cache.Delete(sessionId)
	return nil
}

func (mc *MCache) Exists(sessionId string) (bool, error) {
	return mc.cache.Exists(sessionId), nil
}

func (mc *MCache) Keys(sessionId string) ([]string, error) {
	entry, ok := mc.get(sessionId)
	if !ok || entry == nil {
		return nil, SessionNotExists
	}

	l := len(entry)
	keys := make([]string, l)
	i := 0
	for k, _ := range entry {
		keys[i] = k
		i++
	}
	return keys, nil
}

func (mc *MCache) Init(options string) error {
	return nil
}

func init() {
	Register("default", newMCache())
}
