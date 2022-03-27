logs
[![Go](https://github.com/issue9/logs/workflows/Go/badge.svg)](https://github.com/issue9/logs/actions?query=workflow%3AGo)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/issue9/logs/v4)](https://pkg.go.dev/github.com/issue9/logs/v4)
![Go version](https://img.shields.io/github/go-mod/go-version/issue9/logs)
======

```go
import "github.com/issue9/logs/v4"

l := logs.New()
l.Debug("debug start...")
l.Debugf("%v start...", "debug")
l.DEBUG().Print("debug start...")
```

安装
---

```shell
go get github.com/issue9/logs/v4
```

版权
---

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
