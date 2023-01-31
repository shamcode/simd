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
- [Wrap QueryBuilder](https://github.com/shamcode/simd/blob/master/_examples/wrap-query-builder)


Benchmarks
------
```go
goos: linux
goarch: amd64
pkg: github.com/shamcode/simd/benchmarks
cpu: 11th Gen Intel(R) Core(TM) i7-11700K @ 3.60GHz
Benchmark_SIMDVsSQLite/10_simd-16         	10237388	       592.0 ns/op	     504 B/op	      12 allocs/op
Benchmark_SIMDVsSQLite/10_sqlite-16       	 2908093	      1956 ns/op	     576 B/op	      25 allocs/op
Benchmark_SIMDVsSQLite/100_simd-16        	  420433	     14340 ns/op	   12102 B/op	     288 allocs/op
Benchmark_SIMDVsSQLite/100_sqlite-16      	  124260	     48167 ns/op	   13824 B/op	     600 allocs/op
Benchmark_SIMDVsSQLite/1000_simd-16       	   39694	    153177 ns/op	  125568 B/op	    2988 allocs/op
Benchmark_SIMDVsSQLite/1000_sqlite-16     	   10000	    512700 ns/op	  143424 B/op	    6225 allocs/op
Benchmark_SIMDVsSQLite/5000_simd-16       	    7059	    872707 ns/op	  645894 B/op	   16977 allocs/op
Benchmark_SIMDVsSQLite/5000_sqlite-16     	    2311	   2565130 ns/op	  735330 B/op	   33213 allocs/op
Benchmark_SIMDVsSQLite/10000_simd-16      	    3266	   1852537 ns/op	 1296400 B/op	   34478 allocs/op
Benchmark_SIMDVsSQLite/10000_sqlite-16    	    1018	   5202534 ns/op	 1475333 B/op	   66963 allocs/op
Benchmark_SIMDVsSQLite/50000_simd-16      	     608	   9815205 ns/op	 6496680 B/op	  174480 allocs/op
Benchmark_SIMDVsSQLite/50000_sqlite-16    	     231	  26201932 ns/op	 7395358 B/op	  336963 allocs/op

```

License
-------
Distributed under MIT License, please see license file within the code for more details.