// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package rotate

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/issue9/assert/v4"
)

var _ io.WriteCloser = &rotate{}

func TestNew(t *testing.T) {
	a := assert.New(t, false)

	a.NotError(os.RemoveAll("./testdata"))
	w, err := New("01-02-%i", "./testdata", 100)
	a.NotError(err).NotNil(w)
	a.Equal(w.(*rotate).size, 100)

	loop := 100
	for i := 0; i < loop; i++ {
		// 加个延时，否则全部到一个文件中
		time.Sleep(60 * time.Millisecond)

		size, err := w.Write([]byte("1024\n"))
		a.NotEqual(size, 0).NotError(err)
	}

	files, err := os.ReadDir(w.(*rotate).dir)
	a.NotError(err)
	a.Equal(len(files), int64(loop*len("1024\n"))/w.(*rotate).size)
}
