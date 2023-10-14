logs
[![Go](https://github.com/issue9/logs/actions/workflows/go.yml/badge.svg)](https://github.com/issue9/logs/actions/workflows/go.yml)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/issue9/logs/v6)](https://pkg.go.dev/github.com/issue9/logs/v6)
![Go version](https://img.shields.io/github/go-mod/go-version/issue9/logs)
[![codecov](https://codecov.io/gh/issue9/logs/branch/master/graph/badge.svg)](https://codecov.io/gh/issue9/logs)
======

高性能日志库

```text
goos: darwin
goarch: amd64
pkg: github.com/imkira/go-loggers-bench
cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
BenchmarkLogsTextPositive-4   	100000000	       320.9 ns/op	      40 B/op	       2 allocs/op
BenchmarkLogsTextNegative-4   	1000000000	         9.407 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogsJSONNegative-4   	1000000000	        11.42 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogsJSONPositive-4   	65887180	       578.6 ns/op	      40 B/op	       2 allocs/op
```

```go
import "github.com/issue9/logs/v6"

l := logs.New(logs.NewTextHandler(...))
l.DEBUG().Print("debug start...")

erro := l.With(logs.LevelError, map[string]interface{}{"k1":"v1"})
erro.Printf("带默认参数 k1=v1") // 不用 With 指定 k1，err 全都自动带上此参数
```

安装
---

```shell
go get github.com/issue9/logs/v6
```

版权
---

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
