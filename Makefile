run_test:
	go test -v -race  -covermode=atomic -coverprofile=coverage.out -coverpkg=./... ./... && echo "Tests finished with success"

coverage:
	go tool cover -html=coverage.out

test: run_test coverage

bench:
	go test -bench=. -benchmem ./benchmarks/...