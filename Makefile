run_test:
	go test -v -race  -covermode=atomic -coverprofile=coverage.out -coverpkg=./... ./... && echo "Tests finished with success"

coverage:
	go tool cover -html=coverage.out

test: run_test coverage

bench_query_builder:
	go test -bench=. -benchmem -benchtime=7s -run=^_ ./query/...

bench_comparing:
	go test -bench=. -benchmem -benchtime=5s ./benchmarks/...

make bench: bench_query_builder bench_comparing