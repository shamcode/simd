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

bench_comparing:
	go test -bench=. -benchmem -benchtime=5s ./benchmarks/...

make bench: bench_query_builder bench_comparing