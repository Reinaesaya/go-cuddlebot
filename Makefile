BIN_DIR ?= bin
LINUX_BIN_DIR ?= bin-arm-linux
EXECUTABLES = cuddled cuddlespeak
EXECUTABLES_DEST = $(EXECUTABLES:%=$(BIN_DIR)/%) $(EXECUTABLES:%=$(LINUX_BIN_DIR)/%)

build: $(EXECUTABLES) $(EXECUTABLES_DEST)

clean:
	rm -fr $(BIN_DIR) $(LINUX_BIN_DIR)

$(BIN_DIR):
	mkdir -p $@

$(BIN_DIR)/%: %/main.go $(BIN_DIR)
	go build -o $@ $<

$(LINUX_BIN_DIR):
	mkdir -p $@

$(LINUX_BIN_DIR)/%: %/main.go $(LINUX_BIN_DIR)
	GOARCH=arm GOARM=7 GOOS=linux go build -o $@ $<

.PHONY: build clean
