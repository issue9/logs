// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package rotate 提供一个可以按文件大小进行分割的 io.Writer 实例。
package rotate

import (
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Rotate 可按大小进行分割的文件
//  import "log"
//  // 每个文件以 100M 大小进行分割，以日期名作为文件名保存在 /var/log 下。
//  f,_ := NewRotate("debug-%y%m%i", "/var/log", 100*1024*1024)
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
//
// format 文件名格式，可以包含以下格式内容：
//  %y 表示两位数的年份；
//  %Y 四位数的年份；
//  %m 表示月份；
//  %d 表示当前月中的第几天；
//  %h 表示 12 进制的小时；
//  %H 表示 24 进制的小时；
//  %i 表示同一时间段内的，多个的文件的计数器。
// dir 为文件保存的目录，若不存在会尝试创建。
// size 为每个文件的最大尺寸，单位为 byte。
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

// 打开一个日志文件
func (r *Rotate) open() error {
	if r.w != nil {
		r.w.Close()
	}

	now := time.Now()
	prefix := now.Format(r.prefix)
	suffix := now.Format(r.suffix)

	index, err := getIndex(r.dir, prefix, suffix)
	if err != nil {
		return err
	}

	path := filepath.Join(r.dir, prefix+strconv.Itoa(index)+suffix)
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			goto CREATE
		}
		return err
	}

	if stat.Size() < r.size {
		r.w, err = os.OpenFile(path, os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			return err
		}

		r.wSize = 0
		return nil
	}

	index++

CREATE:
	path = filepath.Join(r.dir, prefix+strconv.Itoa(index)+suffix)
	if r.w, err = os.Create(path); err != nil {
		return err
	}

	r.wSize = 0
	return nil
}

func (r *Rotate) Write(buf []byte) (int, error) {
	if (r.wSize > r.size) || r.w == nil {
		if err := r.open(); err != nil {
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
