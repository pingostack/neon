GO_LDFLAGS = -ldflags "-s -w"
all: package

.phony: init neon package clean

init:
	go mod tidy
	go generate ./...

neon: init
	mkdir -p build/bin
	go build -o build/bin/neon $(GO_LDFLAGS) cmd/neon/main.go

package: neon
	mkdir -p build/config
	cp -n config.yml build/config/config.yml || true

clean:
	rm -rf build

help:
	@echo "make - same as make all"
	@echo "make all - init, build neon, package"
	@echo "make init - init go mod"
	@echo "make neon - build neon"
	@echo "make package - package neon"
	@echo "make clean - clean build"
	@echo "make help - show help"
