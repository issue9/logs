logs
[![Go](https://github.com/issue9/logs/workflows/Go/badge.svg)](https://github.com/issue9/logs/actions?query=workflow%3AGo)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/issue9/logs/v3)](https://pkg.go.dev/github.com/issue9/logs/v3)
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
import "github.com/issue9/logs/v3/config"
import "github.com/issue9/logs/v3"

cfg, _ := config.ParseFile("./logs.xml")
l,err := logs.New(cfg)
l.Debug("debug start...")
l.Debugf("%v start...", "debug")
l.DEBUG().Println("debug start...")
```

安装
---

```shell
go get github.com/issue9/logs/v3
```

版权
---

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
