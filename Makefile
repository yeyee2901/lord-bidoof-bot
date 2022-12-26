# GOLANG VARIABLES
GO_TEST_FLAGS 	+= 	-v -c -coverpkg ./...

# override-able
GO_RUN_TEST		= 	^Test

.PHONY: run test

run:
	go run ./cmd/grpc-controller

test:
	mkdir -p test/config
	go test ${GO_TEST_FLAGS} -o ./test/config/compiled ./pkg/config
	mkdir -p test/telegram
	go test ${GO_TEST_FLAGS} -o ./test/telegram/compiled ./pkg/telegram

test_telegram: test
	./test/telegram/compiled -test.run ${GO_RUN_TEST} -test.count=1 -test.coverprofile=./test/telegram/coverage

test_config: test
	./test/config/compiled -test.run TestLoadConfig -test.count=1 -test.coverprofile=./test/config/coverage
