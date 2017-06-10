coveralls:
	goveralls -service=travis-ci

test:
	go test -v -race ./storage

all: test
