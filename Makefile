.PHONY: build
build:
	go build -mod readonly -o build/simappd ./example/cmd/simappd
	go build -mod readonly -o build/simappcli ./example/cmd/simappcli

.PHONY: protoc
protoc:
	bash ./scripts/protocgen.sh

.PHONY: test
test:
	go test -v ./x/...
