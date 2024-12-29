OUT_DIR=build

.PHONY: build build-program test install docs clean 

build:
	mkdir -p $(OUT_DIR)
	go build -o $(OUT_DIR)/2rm ./src/main.go
	pandoc -s --to man ./README.md -o ./$(OUT_DIR)/2rm.1

build-program:
	mkdir -p $(OUT_DIR)
	go build -o $(OUT_DIR)/2rm ./src/main.go

test:
	go test ./src/...

install:
	scripts/install_artifacts.sh $(OUT_DIR)

docs:
	mkdir -p $(OUT_DIR)
	pandoc -s --to man ./README.md -o ./$(OUT_DIR)/2rm.1

clean:
	rm -rf $(OUT_DIR)
