// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package rotate 提供一个可以按文件大小进行分割的 [io.Writer] 实例
//
//	import "log"
//	// 每个文件以 100M 大小进行分割，以日期名作为文件名保存在 /var/log 下。
//	f,_ := New("debug-200601%i", "/var/log", 100*1024*1024)
//	l := log.New(f, "DEBUG", log.LstdFlags)
package rotate

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type rotate struct {
	dir    string // 文件的保存目录
	size   int64  // 每个文件的最大尺寸
	prefix string
	suffix string

	w     *os.File // 当前正在写的文件
	wSize int64    // 当前正在写的文件大小
}

// New 声明按大小进行文件分割的对象
//
// format 文件名格式，除了标准库支持的时间格式之外，还需要包含
// %i 占位符，表示同一时间段内的产生多个文件时的计数器。比如：
//
//	2006-01-02-15-04-%i-01-02.log
//
// dir 为文件保存的目录，若不存在会尝试创建。
// size 为每个文件的最大尺寸，单位为 byte。
func New(format, dir string, size int64) (io.WriteCloser, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	stat, err := os.Stat(dir)
	if (err != nil && !os.IsExist(err)) || !stat.IsDir() {
		if !errors.Is(err, os.ErrNotExist) {
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

	p, s, err := cutString(format)
	if err != nil {
		return nil, err
	}

	return &rotate{
		dir:    dir,
		prefix: p,
		suffix: s,
		size:   size,
	}, nil
}

// 打开一个日志文件
func (r *rotate) open() error {
	if r.w != nil {
		if err := r.w.Close(); err != nil {
			return err
		}
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

func (r *rotate) Write(buf []byte) (int, error) {
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

func (r *rotate) Close() error {
	if r.w == nil {
		return nil
	}
	return r.w.Close()
}
