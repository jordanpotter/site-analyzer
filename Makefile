default: build

build:
	go build

install:
	go install

lint:
	gometalinter $(shell glide novendor) --deadline 300s

test:
	go test -v $(shell glide novendor)
