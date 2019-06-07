.PHONY: test install cover

TEST_RESULTS ?= ./test-results

all: build

build: test
	go build
	#cd cmd;\
	#go build -o ../capi;

test:
	mkdir -p $(TEST_RESULTS)
#	go get github.com/jstemmer/go-junit-report
	go test ./... -v -vet all | go-junit-report > $(TEST_RESULTS)/report.xml

install:
	# complete  : tool name
	# -C        : this is a command
	# capi      : custom complete tool
	# capi      : name of the target tool (to match on the cli)
	echo complete -C capi capi

tail:
	rm -f complete.log
	touch complete.log
	tail complete.log -f

cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out