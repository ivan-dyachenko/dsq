#
# Makefile
#
BINARY_NAME = dsq
VERSION = snapshot
GHRFLAGS =
.PHONY: buildx release

default: build

build:
	go build -o $(BINARY_NAME) -v ./...

test:
	go test -v ./...

clean:
	go clean
	rm -f $(BINARY_NAME)

run:
	go build -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

buildx:
	goxc -d=pkg -pv=$(VERSION)

release:
	ghr  -u ivan-dyachenko  $(GHRFLAGS) v$(VERSION) pkg/$(VERSION)
