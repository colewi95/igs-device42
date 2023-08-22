HOSTNAME=github.com
NAMESPACE=chopnico
NAME=device42
BINARY=terraform-provider-${NAME}
VERSION=0.3.15
OS_ARCH=linux_amd64

default: install

.PHONY: build

build:
	go fmt ./...
	go build -o /tmp/${BINARY}

release:
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv /tmp/${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
docs:
	go generate

clean:
	rm -rf docs/
