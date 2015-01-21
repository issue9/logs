// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writers

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/issue9/assert"
)

var _ io.WriteCloser = &Rotate{}

// 清空指定目录下的所有内容。
func clearDir(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			if err = os.RemoveAll(dir + file.Name()); err != nil {
				return err
			}
		} else {
			if err = os.Remove(dir + file.Name()); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestRotate(t *testing.T) {
	a := assert.New(t)

	w, err := NewRotate("test_", "./testdata", 100)
	a.NotError(err)
	a.NotNil(w)
	a.Equal(w.size, 100)

	clearDir(w.dir)

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
	a.Equal(len(files), loop*len("1024\n")/w.size)
}
