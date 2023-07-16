logs
[![Go](https://github.com/issue9/logs/workflows/Go/badge.svg)](https://github.com/issue9/logs/actions?query=workflow%3AGo)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/issue9/logs/v5)](https://pkg.go.dev/github.com/issue9/logs/v5)
![Go version](https://img.shields.io/github/go-mod/go-version/issue9/logs)
[![codecov](https://codecov.io/gh/issue9/logs/branch/master/graph/badge.svg)](https://codecov.io/gh/issue9/logs)
======

高性能日志库

```
BenchmarkLogsTextPositive     	19077283	       315.7 ns/op	      48 B/op	       1 allocs/op
BenchmarkLogsTextPositive-2   	32050080	       196.7 ns/op	      48 B/op	       1 allocs/op
BenchmarkLogsTextPositive-4   	31083766	       177.2 ns/op	      48 B/op	       1 allocs/op
BenchmarkLogsTextNegative     	233248940	        24.56 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogsTextNegative-2   	440099184	        13.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogsTextNegative-4   	668924367	         9.313 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogsJSONNegative     	80700519	        70.91 ns/op	      16 B/op	       1 allocs/op
BenchmarkLogsJSONNegative-2   	159696355	        41.04 ns/op	      16 B/op	       1 allocs/op
BenchmarkLogsJSONNegative-4   	114847047	        45.25 ns/op	      16 B/op	       1 allocs/op
BenchmarkLogsJSONPositive     	 6550818	       886.1 ns/op	      64 B/op	       2 allocs/op
BenchmarkLogsJSONPositive-2   	11191477	       517.9 ns/op	      64 B/op	       2 allocs/op
BenchmarkLogsJSONPositive-4   	11880020	       515.5 ns/op	      64 B/op	       2 allocs/op
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
