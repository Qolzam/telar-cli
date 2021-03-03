Version := $(shell git describe --tags --dirty)
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-s -w -X github.com/Qolzam/telar-cli/cmd.Version=$(Version) -X github.com/Qolzam/telar-cli/cmd.GitCommit=$(GitCommit)"
PLATFORM := $(shell ./hack/platform-tag.sh)
SOURCE_DIRS = cmd pkg main.go
export GO111MODULE=on

.PHONY: all
all: clean gofmt test build dist hash

.PHONY: clean
clean:
	@rm -rf ./ui/build
	@rm -rf ./pkged.go

.PHONY: build
build:
	go build

.PHONY: gofmt
gofmt:
	@test -z $(shell gofmt -l -s $(SOURCE_DIRS) ./ | tee /dev/stderr) || (echo "[WARN] Fix formatting issues with 'make fmt'" && exit 1)

.PHONY: test
test:
	CGO_ENABLED=0 go test $(shell go list ./... | grep -v /vendor/|xargs echo) -cover

.PHONY: dist
dist:
	@cd ./ui && yarn install
	@cd ./ui && yarn build
	pkger -include /ui/build
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/telar-cli
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/telar-cli-darwin
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/telar-cli-armhf
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/telar-cli-arm64
	CGO_ENABLED=0 GOOS=windows go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/telar-cli.exe

.PHONY: hash
hash:
	rm -rf bin/*.sha256 && ./hack/hashgen.sh
