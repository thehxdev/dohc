GO := go
BIN := dohc
CROSS_PLAT_BUILD_DIR := ./build

$(BIN): $(wildcard *.go) $(wildcard */*.go)
	CGO_ENABLED=0 $(GO) build -ldflags='-s -w -buildid=' .

cross-plat:
	@mkdir -p $(CROSS_PLAT_BUILD_DIR)
	@python3 build.py

clean:
	rm -rf $(BIN) result*.txt $(CROSS_PLAT_BUILD_DIR)
	go clean
