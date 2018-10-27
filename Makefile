.PHONY: build test install

all: build

build:
	cd cmd;\
	go build -vet all -o ../capi;

test:
	go test ./... -vet all

install:
	complete -C ./capi capi


