.PHONY: build

# TODO: change PATH to include bin for test
# TODO: support release on windows (zip not on windows)

format:
	cd src/go; gofmt -w .

build: format
	rm -rf bin; \
	mkdir -p bin; \
	cd src/go; \
	go build -o ../../bin ./...  
	cp src/bash/* bin

test: build
	cd src/go; go test ./...

release: build
ifndef RELEASE_VERSION
	$(error RELEASE_VERSION is not set)
endif
ifndef PLATFORM
	$(error PLATFORM is not set)
endif
	rm -rf build/zip
	mkdir -p build/zip/stacked-diff-workflow
	cp bin/* build/zip/stacked-diff-workflow
	cd build/zip; \
	zip -vr stacked-diff-workflow-${PLATFORM}-$(RELEASE_VERSION).zip stacked-diff-workflow/ -x "*.DS_Store"
