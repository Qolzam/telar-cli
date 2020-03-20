

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
	# @go build -o ./telar-cli
	env GOOS=windows GOARCH=amd64 go build -o  ./telar-cli.exe
	@echo "[✔️] Build complete!"

.PHONY: run

run:
	# @open ./telar-cli
	@echo "[✔️] App is running!"
