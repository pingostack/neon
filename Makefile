GO_LDFLAGS = -ldflags "-s -w"
all: package

.phony: init build_neon package clean

init:
	go mod tidy
	go generate ./...

build_neon: init
	mkdir -p build/bin
	go build -o build/bin/neon $(GO_LDFLAGS) cmd/neon/main.go

package: build_neon
	mkdir -p build/config
	cp -n config.yml build/config/config.yml || true

clean:
	rm -rf build
