

.PHONY: all

all: clean build run

.PHONY: clean

clean:
	@rm -rf ./ui/build
	@echo "[✔️] Clean complete!"

.PHONY: build

build:
	@cd ./ui && yarn install
	@cd ./ui && yarn build
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o bin/telar
	GO111MODULE=on CGO_ENABLED=0 GOOS=darwin go build -o bin/telar-darwin
	@echo "[✔️] Build complete!"
	packr clean

.PHONY: run

run:
	@open ./bin/telar-darwin
	@echo "[✔️] App is running!"
