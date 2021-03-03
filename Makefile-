

.PHONY: all

all: clean build run

.PHONY: clean

clean:
	@rm -rf ./ui/build
	@rm -rf ./pkged.go
	@echo "[✔️] Clean complete!"

.PHONY: build

build:
	@cd ./ui && yarn install
	@cd ./ui && yarn build
	pkger -include /ui/build
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o bin/telar
	GO111MODULE=on CGO_ENABLED=0 GOOS=darwin go build -o bin/telar-darwin
	@echo "[✔️] Build complete!"
	
.PHONY: run

run:
	@echo "[✔️] App is running!"
