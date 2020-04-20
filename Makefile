FILES?=main.go
PLATFORM?=darwin
ARCHITECTURE?=amd64

BUILDTIME=`date "+%F %T%Z"`
VERSION=`git describe --tags`

build:
	GOOS=$(PLATFORM) GOARCH=$(ARCHITECTURE) go build -ldflags="-X 'main.buildTime=$(BUILDTIME)' -X 'main.version=$(VERSION)' -s -w -extldflags '-static'" -o bin/erply-middleware $(FILES)

run:
	go run $(FILES)

clean:
	rm -rf bin

test:
	go test ./...

cover:
	go test ./... -coverprofile cover.out
