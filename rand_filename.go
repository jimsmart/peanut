// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package peanut

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// This is copied from https://golang.org/src/io/ioutil/tempfile.go

// Random number state.
// We generate random temporary file names so that there's a good
// chance the file doesn't exist yet - keeps the number of tries in
// TempFile to a minimum.
var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

func nextRandom() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

// This began life as os.TempFile, but has been refactored somewhat from the original.
func randomTempFilename(prefix, suffix string) (string, error) {
	dir := os.TempDir()
	nconflict := 0
	for i := 0; i < 10000; i++ {
		name := filepath.Join(dir, prefix+nextRandom()+suffix)
		_, err := os.Stat(name)
		if !os.IsExist(err) {
			return name, nil
		}
		if nconflict++; nconflict > 10 {
			randmu.Lock()
			rand = reseed()
			randmu.Unlock()
		}
	}
	return "", errors.New("Random filename creation failed (too many conflicts)")
}
