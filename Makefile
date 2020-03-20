

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
	env GOOS=windows GOARCH=amd64 go build -o  ./telar-cli.exe
	# @go build -o ./telar-cli
	@echo "[✔️] Build complete!"

.PHONY: run

run:
	# @open ./telar-cli
	@echo "[✔️] App is running!"
