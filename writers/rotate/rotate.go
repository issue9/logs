// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package rotate

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Rotate 可按大小进行分割的文件
//  import "log"
//  // 每个文件以 100M 大小进行分割，以日期名作为文件名保存在 /var/log 下。
//  f,_ := NewRotate("debug-", "/var/log", 100*1024*1024)
//  l := log.New(f, "DEBUG", log.LstdFlags)
type Rotate struct {
	dir    string // 文件的保存目录
	size   int64  // 每个文件的最大尺寸
	prefix string
	suffix string

	w     *os.File // 当前正在写的文件
	wSize int64    // 当前正在写的文件大小
}

// New 新建 Rotate。
// prefix 文件名格式。
// dir 为文件保存的目录，若不存在会尝试创建。
// size 为每个文件的最大尺寸，单位为 byte。size 应该足够大，如果 size
// 的大小不足够支撑一秒钟产生的量，则会继续在原有文件之后追加内容。
func New(format, dir string, size int64) (*Rotate, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	stat, err := os.Stat(dir)
	if (err != nil && !os.IsExist(err)) || !stat.IsDir() {
		if !os.IsNotExist(err) {
			return nil, err
		}

		// 尝试创建目录
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}

		// 创建目录成功，重新获取状态
		if _, err = os.Stat(dir); err != nil {
			return nil, err
		}
	}

	p, s, err := parseFormat(format)
	if err != nil {
		return nil, err
	}

	return &Rotate{
		dir:    dir,
		prefix: p,
		suffix: s,
		size:   size,
	}, nil
}

// 初始化一个新的文件对象
func (r *Rotate) init() error {
	if r.w != nil {
		r.w.Close()
	}

	now := time.Now()
	prefix := now.Format(r.prefix)
	suffix := now.Format(r.suffix)
	fs, err := ioutil.ReadDir(r.dir)
	if err != nil {
		return err
	}

	var index int
	var path string
	var stat os.FileInfo
	for ; index < len(fs); index++ {
		name := prefix + strconv.Itoa(index) + suffix
		path = filepath.Join(r.dir, name)

		stat, err = os.Stat(path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	if stat != nil && stat.Size() < r.size {
		r.w, err = os.Open(path)
		if err != nil {
			return err
		}

		r.wSize = 0
		return nil
	}

	name := prefix + strconv.Itoa(index) + suffix
	path = filepath.Join(r.dir, name)
	if r.w, err = os.Create(path); err != nil {
		return err
	}

	r.wSize = 0
	return nil
}

func (r *Rotate) Write(buf []byte) (int, error) {
	if (r.wSize > r.size) || r.w == nil {
		if err := r.init(); err != nil {
			return 0, err
		}
	}

	size, err := r.w.Write(buf)
	if err != nil {
		return 0, err
	}

	r.wSize += int64(size)

	return size, nil
}

// Close 关闭文件
func (r *Rotate) Close() error {
	if r.w == nil {
		return nil
	}

	return r.w.Close()
}

// Flush 实现接口 Flusher.Flush()
func (r *Rotate) Flush() {
	r.Close()
}
