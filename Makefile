.PHONY: prepare build test clean

prepare:
ifeq (,$(wildcard bin/build))
	$(info ******************** downloading build script ********************)
	./scripts/fetch-build-script.sh
endif

build: prepare
	$(info ******************** building binary ********************)
	./bin/build

test:
	$(info ******************** running tests ********************)
	go test -v ./...

clean:
	./bin/build --clean
