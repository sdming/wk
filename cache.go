// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"time"
)

// Cacher is interface that wraps cache methods
type Cacher interface {

	// Get return cached data by key
	Get(key string) ([]byte, error)

	// Set add data to cache or update cached data by key
	Set(key string, value interface{}, expiration time.Duration) error

	// Remove remove cached data from cache
	Remove(key string) error
}
