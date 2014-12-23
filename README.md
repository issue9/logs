logs [![Build Status](https://travis-ci.org/issue9/logs.svg?branch=master)](https://travis-ci.org/issue9/logs)
======

一个可配置的日志服务包。可以通过xml自定义日志输出：
```xml
<?xml version="1.0" encoding="utf-8" ?>
<logs>
    <debug>
        <buffer size="10">
            <rotate dir="/var/log/" size="5M" />
            <stmp username=".." password=".." />
        </buffer>
    </debug>
    <info>
        ....
    </info>
</logs>
```


```go
logs.InitFromXmlFile("./config.xml")// 用xml初始化logs
logs.Debug("debug start...")
logs.Debugf("%v start...", "debug")
logs.DEBUG.Println("debug start...")
```

### 安装

```shell
go get github.com/issue9/logs
```


### 文档

[![Go Walker](http://gowalker.org/api/v1/badge)](http://gowalker.org/github.com/issue9/logs)
[![GoDoc](https://godoc.org/github.com/issue9/logs?status.svg)](https://godoc.org/github.com/issue9/logs)


### 版权

[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/issue9/logs/blob/master/LICENSE)
