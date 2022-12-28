# PROTOBUF
PROTO_REPO = github.com/yeyee2901/proto-lord-bidoof-bot@master

# GOLANG VARIABLES
GO_TEST_FLAGS 	+= 	-v -c -coverpkg ./...

# override-able
GO_RUN_TEST		= 	^Test

.PHONY: run test

run:
	go run ./cmd/grpc-controller

update_proto:
	go get -x -u ${PROTO_REPO}

test:
	mkdir -p test/config
	go test ${GO_TEST_FLAGS} -o ./test/config/compiled ./pkg/config
	mkdir -p test/telegram
	go test ${GO_TEST_FLAGS} -o ./test/telegram/compiled ./pkg/telegram
	mkdir -p test/datasource
	go test ${GO_TEST_FLAGS} -o ./test/datasource/compiled ./pkg/datasource

test_telegram: test
	./test/telegram/compiled -test.v -test.run ${GO_RUN_TEST} -test.count=1 -test.coverprofile=./test/telegram/coverage

test_config: test
	./test/config/compiled -test.v -test.run TestLoadConfig -test.count=1 -test.coverprofile=./test/config/coverage

test_db: test
	./test/datasource/compiled -test.v test.run TestGetPrivateChatWithQueryFilter -test.count=1 -test.coverprofile=./test/datasource/db-coverage
