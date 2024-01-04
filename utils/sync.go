package utils

import "sync"

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.zenda file.

func OnceFunc(f func()) func() {
	var (
		once  sync.Once
		valid bool
		p     any
	)
	// Construct the inner closure just once to reduce costs on the fast path.
	g := func() {
		defer func() {
			p = recover()
			if !valid {
				// Re-panic immediately so on the first call the user gets a
				// complete stack trace into f.
				panic(p)
			}
		}()
		f()
		valid = true // Set only if f does not panic
	}
	return func() {
		once.Do(g)
		if !valid {
			panic(p)
		}
	}
}
