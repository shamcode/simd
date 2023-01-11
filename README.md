simd
=================
![Project status](https://img.shields.io/badge/version-0.0.3-green.svg)
![Build](https://github.com/shamcode/simd/actions/workflows/workflow.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/shamcode/simd)](https://goreportcard.com/report/github.com/shamcode/simd)
[![Coverage Status](https://coveralls.io/repos/github/shamcode/simd/badge.svg?branch=master)](https://coveralls.io/github/shamcode/simd?branch=master)
[![GoDoc](https://godoc.org/github.com/shamcode/simd?status.svg)](https://pkg.go.dev/github.com/shamcode/simd)
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
Benchmark_CompareSIMDWithSQLite/10_simd-16         6106621	      1035 ns/op	     504 B/op	      12 allocs/op
Benchmark_CompareSIMDWithSQLite/10_sqlite-16       2933868	      2088 ns/op	     576 B/op	      25 allocs/op
Benchmark_CompareSIMDWithSQLite/100_simd-16         236976	     25478 ns/op	   12102 B/op	     288 allocs/op
Benchmark_CompareSIMDWithSQLite/100_sqlite-16       120069	     48522 ns/op	   13824 B/op	     600 allocs/op
Benchmark_CompareSIMDWithSQLite/1000_simd-16         18895	    354572 ns/op	  125571 B/op	    2988 allocs/op
Benchmark_CompareSIMDWithSQLite/1000_sqlite-16       11125	    529965 ns/op	  143424 B/op	    6225 allocs/op
Benchmark_CompareSIMDWithSQLite/5000_simd-16          4045	   2027158 ns/op	  638062 B/op	   15983 allocs/op
Benchmark_CompareSIMDWithSQLite/5000_sqlite-16        2059	   2767068 ns/op	  735330 B/op	   33213 allocs/op
Benchmark_CompareSIMDWithSQLite/10000_simd-16         2557	   2573430 ns/op	 1278173 B/op	   32233 allocs/op
Benchmark_CompareSIMDWithSQLite/10000_sqlite-16       1074	   5329826 ns/op	 1475332 B/op	   66963 allocs/op

```

License
-------
Distributed under MIT License, please see license file within the code for more details.