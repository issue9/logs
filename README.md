logs
[![Build Status](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fissue9%2Flogs%2Fbadge%3Fref%3Dmaster&style=flat)](https://actions-badge.atrox.dev/issue9/logs/goto?ref=master)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/issue9/logs/v2)
======

一个可配置的日志服务包。可以通过 XML 自定义日志输出：

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
logs.InitFromXMLFile("./config.xml") // 用 XML 初始化 logs
logs.Debug("debug start...")
logs.Debugf("%v start...", "debug")
logs.DEBUG().Println("debug start...")
```

安装
---

```shell
go get github.com/issue9/logs
```

版权
---

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
