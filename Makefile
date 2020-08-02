.PHONY: protoc
protoc:
	bash ./scripts/protocgen.sh

.PHONY: test
test:
	go test -v ./x/...