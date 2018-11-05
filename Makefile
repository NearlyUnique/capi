.PHONY: build test install

all: build

build:
	cd cmd;\
	go build -o ../capi;

test:
	go test ./... -vet all

install:
	echo complete -C ./capi capi

tail:
	rm -f complete.log
	touch complete.log
	tail complete.log -f

