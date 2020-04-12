

.PHONY: all

all: clean build run

.PHONY: clean

clean:
	@rm -rf ./ui/build
	@echo "[✔️] Clean complete!"

.PHONY: build

build:
	# @cd ./ui && npm install
	@cd ./ui && npm run build
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/telar
	GO111MODULE=on CGO_ENABLED=0 GOOS=darwin go build -a -installsuffix cgo -o bin/telar-darwin
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -a -installsuffix cgo -o bin/telar-armhf
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -o bin/telar-arm64
	GO111MODULE=on CGO_ENABLED=0 GOOS=windows go build -a -installsuffix cgo -o bin/telar.exe

	@echo "[✔️] Build complete!"

.PHONY: run

run:
	# @open ./telar-cli
	@echo "[✔️] App is running!"
