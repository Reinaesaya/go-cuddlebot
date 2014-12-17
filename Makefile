BIN_DIR ?= bin-arm-linux
EXECUTABLES = cuddled cuddlespeak
EXECUTABLES_DEST = $(EXECUTABLES:%=$(BIN_DIR)/%)

build: $(EXECUTABLES_DEST)

clean:
	rm $(EXECUTABLES_DEST)

$(BIN_DIR):
	mkdir -p $<

$(BIN_DIR)/%: %/main.go $(BIN_DIR)
	GOARCH=arm GOARM=7 GOOS=linux go build -o $@ $<

.PHONY: build clean
