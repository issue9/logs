// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writer

import (
	"fmt"
	"io"
	"os"
	"time"
)

var (
	// 默认的文件权限
	defaultMode os.FileMode = 0660

	// linux下需加上O_WRONLY或是O_RDWR
	defaultFlag = os.O_APPEND | os.O_CREATE | os.O_WRONLY

	// 默认的日志后缀名
	defaultExt = ".log"
)

// 可按大小进行分割的文件
//  import "log"
//  // 每个文件以100M大小进行分割，以日期名作为文件名保存在/var/log下。
//  f,_ := NewRotate("/var/log", 100*1024*1024)
//  l := log.New(f, "DEBUG", log.LstdFlags)
type Rotate struct {
	dir      string // 文件的保存目录
	size     int    // 每个文件的最大尺寸
	basePath string

	w     io.Writer // 当前正在写的文件
	wSize int       // 当前正在写的文件大小
}

// 新建Rotate。
// prefix 文件名前缀。
// dir为文件保存的目录，若不存在会尝试创建。
// size为每个文件的最大尺寸，单位为byte。size应该足够大，如果size
// 的大小不足够支撑一秒钟产生的量，则会继续在原有文件之后追加内容。
func NewRotate(prefix, dir string, size int) (*Rotate, error) {
	info, err := os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		// 尝试创建目录
		if err := os.MkdirAll(dir, defaultMode); err != nil {
			return nil, err
		}

		// 创建目录成功，重新获取状态
		if info, err = os.Stat(dir); err != nil {
			return nil, err
		}
	}

	if !info.Mode().IsDir() {
		return nil, fmt.Errorf("[%v]不是一个目录", dir)
	}

	dir = dir + string(os.PathSeparator)
	return &Rotate{
		dir:      dir,
		basePath: dir + prefix,
		size:     size,
	}, nil
}

// 初始化一个新的文件对象
func (r *Rotate) init() error {
	if r.w != nil {
		r.w.(*os.File).Close()
	}

	name := r.basePath + time.Now().Format("20060102150405") + defaultExt

	var err error
	if r.w, err = os.OpenFile(name, defaultFlag, defaultMode); err != nil {
		return err
	}

	r.wSize = 0

	return nil
}

// io.WriteCloser.Write()
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

// io.WriteCloser.Close()
func (r *Rotate) Close() error {
	if r.w == nil {
		return nil
	}

	return r.w.(*os.File).Close()
}

// Flusher.Flush()
func (r *Rotate) Flush() {
	r.Close()
}
