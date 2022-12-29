simd
=================
![Project status](https://img.shields.io/badge/version-0.0.2-green.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/shamcode/simd)](https://goreportcard.com/report/github.com/shamcode/simd)
[![Coverage Status](https://coveralls.io/repos/github/shamcode/simd/badge.svg?branch=master)](https://coveralls.io/github/shamcode/simd?branch=master)
[![GoDoc](https://godoc.org/github.com/shamcode/simd?status.svg)](https://pkg.go.dev/github.com/shamcode/simd/v0)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

simd (**S**imple **I**n **M**emory **D**atabase) &mdash; is an embeddable golang database with support for conditional queries, custom sorting and custom field types.


Installation
------------
 
Use go get.

    go get github.com/shamcode/simd

Usage
------

##### Examples:

- [Simple](https://github.com/shamcode/simd/blob/master/_examples/common/main.go)
- [Custom Field Type](https://github.com/shamcode/simd/blob/master/_examples/custom-field-time)


Benchmarks
------
```go
goos: linux
goarch: amd64
pkg: github.com/shamcode/simd/benchmarks
cpu: 11th Gen Intel(R) Core(TM) i7-11700K @ 3.60GHz
Benchmark_CompareSIMDWithSQLite/10_simd-16       1284638	      1113 ns/op	     520 B/op	      13 allocs/op
Benchmark_CompareSIMDWithSQLite/10_sqlite-16      487112	      2096 ns/op	     576 B/op	      25 allocs/op
Benchmark_CompareSIMDWithSQLite/100_simd-16        50193	     23755 ns/op	   12486 B/op	     312 allocs/op
Benchmark_CompareSIMDWithSQLite/100_sqlite-16      23046	     50004 ns/op	   13824 B/op	     600 allocs/op
Benchmark_CompareSIMDWithSQLite/1000_simd-16        7252	    302914 ns/op	  129557 B/op	    3237 allocs/op
Benchmark_CompareSIMDWithSQLite/1000_sqlite-16      2053	    524912 ns/op	  143424 B/op	    6225 allocs/op
Benchmark_CompareSIMDWithSQLite/5000_simd-16         792	   2229240 ns/op	  658064 B/op	   17232 allocs/op
Benchmark_CompareSIMDWithSQLite/5000_sqlite-16       404	   2794524 ns/op	  735330 B/op	   33213 allocs/op
Benchmark_CompareSIMDWithSQLite/10000_simd-16        633	   2363879 ns/op	 1318183 B/op	   34733 allocs/op
Benchmark_CompareSIMDWithSQLite/10000_sqlite-16      230	   5318200 ns/op	 1475333 B/op	   66963 allocs/op

```

License
-------
Distributed under MIT License, please see license file within the code for more details.