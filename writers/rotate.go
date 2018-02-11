// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writers

import (
	"os"
	"path/filepath"
	"time"
)

var (
	// 默认的文件权限
	defaultMode os.FileMode = os.ModePerm

	// linux 下需加上 O_WRONLY 或是 O_RDWR
	defaultFlag = os.O_APPEND | os.O_CREATE | os.O_WRONLY

	// 默认的日志后缀名
	defaultExt = ".log"

	filenameFormat = "20060102150405"
)

// Rotate 可按大小进行分割的文件
//  import "log"
//  // 每个文件以 100M 大小进行分割，以日期名作为文件名保存在 /var/log 下。
//  f,_ := NewRotate("/var/log", 100*1024*1024)
//  l := log.New(f, "DEBUG", log.LstdFlags)
type Rotate struct {
	dir      string // 文件的保存目录
	size     int    // 每个文件的最大尺寸
	basePath string

	w     *os.File // 当前正在写的文件
	wSize int      // 当前正在写的文件大小
}

// NewRotate 新建 Rotate。
// prefix 文件名前缀。
// dir 为文件保存的目录，若不存在会尝试创建。
// size 为每个文件的最大尺寸，单位为 byte。size 应该足够大，如果 size
// 的大小不足够支撑一秒钟产生的量，则会继续在原有文件之后追加内容。
func NewRotate(prefix, dir string, size int) (*Rotate, error) {
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
		if err := os.MkdirAll(dir, defaultMode); err != nil {
			return nil, err
		}

		// 创建目录成功，重新获取状态
		if _, err = os.Stat(dir); err != nil {
			return nil, err
		}
	}

	return &Rotate{
		dir:      dir,
		basePath: filepath.Join(dir, prefix),
		size:     size,
	}, nil
}

// 初始化一个新的文件对象
func (r *Rotate) init() (err error) {
	if r.w != nil {
		r.w.Close()
	}

	name := r.basePath + time.Now().Format(filenameFormat) + defaultExt
	if r.w, err = os.OpenFile(name, defaultFlag, defaultMode); err != nil {
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

	r.wSize += size

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
