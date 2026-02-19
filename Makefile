example_common:
	go run ./examples/common -debug

example_custom_field_time:
	go run ./examples/custom-field-time -debug

example_wrap_query_builder:
	go run ./examples/wrap-query-builder -debug

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

lint:
	docker run --rm -t -v ./:/app -v ~/.cache/golangci-lint/v2.10.1:/root/.cache -w /app golangci/golangci-lint:v2.10.1 golangci-lint run

lint_fix:
	docker run --rm -t -v ./:/app -v ~/.cache/golangci-lint/v2.10.1:/root/.cache -w /app golangci/golangci-lint:v2.10.1 golangci-lint run --fix