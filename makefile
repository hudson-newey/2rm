OUT_DIR=build
ENTRY_POINT=./src/main.go
BRAND_NAME=2rm

.PHONY: build build-program test install docs clean 

build:
	mkdir -p $(OUT_DIR)
	go build -o $(OUT_DIR)/$(BRAND_NAME) $(ENTRY_POINT)
	GOOS=windows GOARCH=386 go build -o $(OUT_DIR)/$(BRAND_NAME).exe $(ENTRY_POINT)
	pandoc -s --to man ./README.md -o ./$(OUT_DIR)/$(BRAND_NAME).1

build-program:
	mkdir -p $(OUT_DIR)
	go build -o $(OUT_DIR)/$(BRAND_NAME) ./src/main.go

test:
	go test ./src/...

install:
	scripts/install_artifacts.sh $(OUT_DIR)

docs:
	mkdir -p $(OUT_DIR)
	pandoc -s --to man ./README.md -o ./$(OUT_DIR)/$(BRAND_NAME).1

clean:
	rm -rf $(OUT_DIR)
