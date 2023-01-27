example_common:
	go run ./_examples/common -debug

example_custom_field_time:
	go run ./_examples/custom-field-time -debug

example_wrap_query_builder:
	go run ./_examples/wrap-query-builder -debug

example: example_common example_custom_field_time example_wrap_query_builder

run_test:
	go test -v -race  -covermode=atomic -coverprofile=coverage.out -coverpkg=./... ./... && echo "Tests finished with success"

coverage:
	go tool cover -html=coverage.out

test: example run_test coverage

bench_query_builder:
	go test -bench=. -benchmem -benchtime=7s -run=^_ ./query/...

bench_comparing_sqlite:
	go test -bench=Benchmark_SIMDVsSQLite -benchmem -benchtime=5s -run=^_ ./benchmarks/...

bench_query:
	go test -bench=Benchmark_Query -benchmem -benchtime=5s -run=^_ ./benchmarks/...

bench_set:
	go test -bench=. -benchmem -benchtime=3s -run=^_ ./set/...

bench_concurrent:
	go test -bench=Benchmark_Concurrent -benchmem -benchtime=2000x -run=^_ ./benchmarks/...

bench_indexes:
	go test -bench=Benchmark_Indexes -benchmem -benchtime=1s -run=^_ ./benchmarks/...

bench_indexes_btree:
	go test -bench=Benchmark_BTreeIndexesMaxChildren -benchmem -benchtime=1s -run=^_ ./benchmarks/...

make bench: bench_set bench_query_builder bench_query bench_indexes bench_indexes_btree bench_concurrent bench_comparing_sqlite