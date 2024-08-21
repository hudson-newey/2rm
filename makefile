OUT_DIR=build

build:
	mkdir -p $(OUT_DIR)
	go build -o $(OUT_DIR)/2rm ./src/main.go
