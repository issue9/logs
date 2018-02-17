// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package rotate

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/issue9/assert"
)

var _ io.WriteCloser = &Rotate{}

func TestRotate(t *testing.T) {
	a := assert.New(t)

	a.NotError(os.RemoveAll("./testdata"))
	w, err := New("%i", "./testdata", 100)
	a.NotError(err).NotNil(w)
	a.Equal(w.size, 100)

	loop := 100
	for i := 0; i < loop; i++ {
		// 加个延时，否则全部到一个文件中
		time.Sleep(60 * time.Millisecond)

		size, err := w.Write([]byte("1024\n"))
		a.NotEqual(size, 0)
		a.NotError(err)
	}

	files, err := ioutil.ReadDir(w.dir)
	a.NotError(err)
	a.Equal(len(files), int64(loop*len("1024\n"))/w.size)
}
