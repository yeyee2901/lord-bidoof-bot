# Golang test flags
GO_TEST_FLAGS += -v -c -coverpkg ./...

.PHONY: run test

run:
	go run ./cmd/grpc-controller

test:
	mkdir -p test/config
	go test ${GO_TEST_FLAGS} -o ./test/config/compiled ./pkg/config

test_config: test
	./test/config/compiled -test.run TestLoadConfig -test.count=1 -test.coverprofile=./test/config/coverage
