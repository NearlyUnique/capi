.PHONY: build test install

all: build

build:
	cd cmd;\
	go build -o ../capi;

test:
	go test ./... -vet all

install:
	complete -C ./capi capi


