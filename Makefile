.PHONY: test install cover

all: build

build: test
	go build
	#cd cmd;\
	#go build -o ../capi;

test:
	go test ./... -vet all

install:
	echo complete -C ./capi capi

tail:
	rm -f complete.log
	touch complete.log
	tail complete.log -f

cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out