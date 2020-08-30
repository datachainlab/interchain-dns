.PHONY: build
build:
	go build -mod readonly -o build/simd ./simapp/simd

.PHONY: protoc
protoc:
	bash ./scripts/protocgen.sh

.PHONY: test
test:
	go test -v ./x/...
	go test ./simapp/...