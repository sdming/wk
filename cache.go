// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

type Cacher interface {

	// Get cached item from cache
	Get(key string) ([]byte, error)

	// add or update cache item, expiration in ms
	Set(key string, value interface{}, expiration int) error

	// remove item from cache
	Remove(key string) error
}
