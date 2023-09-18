logs
[![Go](https://github.com/issue9/logs/workflows/Go/badge.svg)](https://github.com/issue9/logs/actions?query=workflow%3AGo)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/issue9/logs/v5)](https://pkg.go.dev/github.com/issue9/logs/v5)
![Go version](https://img.shields.io/github/go-mod/go-version/issue9/logs)
[![codecov](https://codecov.io/gh/issue9/logs/branch/master/graph/badge.svg)](https://codecov.io/gh/issue9/logs)
======

高性能日志库

```
cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
BenchmarkLogsTextPositive     	70519216	       178.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogsTextPositive-2   	100000000	       120.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogsTextPositive-4   	81085801	       141.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogsTextNegative     	484125342	        24.35 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogsTextNegative-2   	1000000000	        14.10 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogsTextNegative-4   	1000000000	         9.033 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogsJSONNegative     	461385842	        25.81 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogsJSONNegative-2   	886811012	        13.65 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogsJSONNegative-4   	1000000000	        10.70 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogsJSONPositive     	20922458	       580.0 ns/op	      40 B/op	       2 allocs/op
BenchmarkLogsJSONPositive-2   	37168568	       318.2 ns/op	      40 B/op	       2 allocs/op
BenchmarkLogsJSONPositive-4   	36571064	       315.3 ns/op	      40 B/op	       2 allocs/op
```

```go
import "github.com/issue9/logs/v5"

l := logs.New(logs.NewTextHandler(...))
l.DEBUG().Print("debug start...")

erro := l.With(logs.LevelError, map[string]interface{}{"k1":"v1"})
erro.Printf("带默认参数 k1=v1") // 不用 With 指定 k1，err 全都自动带上此参数
```

安装
---

```shell
go get github.com/issue9/logs/v5
```

版权
---

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
