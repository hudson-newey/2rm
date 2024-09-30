OUT_DIR=build

.PHONY: clean build

build:
	mkdir -p $(OUT_DIR)
	go build -o $(OUT_DIR)/2rm ./src/main.go

test:
	go test ./src/...

clean:
	rm -rf $(OUT_DIR)

install:
	cp ./build/2rm ~/.local/bin/2rm
