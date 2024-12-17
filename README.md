simd
=================
![Project status](https://img.shields.io/badge/version-0.1.2-green.svg)
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
```text
goos: linux
goarch: amd64
pkg: github.com/shamcode/simd/benchmarks
cpu: 11th Gen Intel(R) Core(TM) i7-11700K @ 3.60GHz
Benchmark_SIMDVsSQLite/10_simd-16         	 7684887	      1003 ns/op	     528 B/op	      13 allocs/op
Benchmark_SIMDVsSQLite/10_sqlite-16       	 3042801	      1951 ns/op	     528 B/op	      16 allocs/op
Benchmark_SIMDVsSQLite/100_simd-16        	  287352	     24983 ns/op	   12680 B/op	     312 allocs/op
Benchmark_SIMDVsSQLite/100_sqlite-16      	  128025	     47585 ns/op	   12672 B/op	     384 allocs/op
Benchmark_SIMDVsSQLite/1000_simd-16       	   23180	    269121 ns/op	  131564 B/op	    3237 allocs/op
Benchmark_SIMDVsSQLite/1000_sqlite-16     	   10000	    503999 ns/op	  131472 B/op	    3984 allocs/op
Benchmark_SIMDVsSQLite/5000_simd-16       	    4161	   1710954 ns/op	  676039 B/op	   18226 allocs/op
Benchmark_SIMDVsSQLite/5000_sqlite-16     	    2257	   2581590 ns/op	  675379 B/op	   21972 allocs/op
Benchmark_SIMDVsSQLite/10000_simd-16      	    2406	   3341910 ns/op	 1356782 B/op	   36977 allocs/op
Benchmark_SIMDVsSQLite/10000_sqlite-16    	    1135	   5211279 ns/op	 1355384 B/op	   44472 allocs/op
Benchmark_SIMDVsSQLite/50000_simd-16      	     567	  10497232 ns/op	 6797117 B/op	  186979 allocs/op
Benchmark_SIMDVsSQLite/50000_sqlite-16    	     224	  26435995 ns/op	 6795416 B/op	  224472 allocs/op

```

License
-------
Distributed under MIT License, please see license file within the code for more details.