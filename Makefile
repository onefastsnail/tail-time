.PHONY: *

lint:
	golangci-lint run

test:
	go test -cover ./...

cov-func:
	mkdir -p build
	FILE=`mktemp build/coverage.XXXX` && \
		go test -coverprofile=$${FILE} ./... && \
		go tool cover -func=$${FILE} && \
		rm $${FILE}

cov-html:
	mkdir -p build
	FILE=`mktemp build/coverage.XXXX` && \
		go test -coverprofile=$${FILE} ./... && \
		go tool cover -html=$${FILE} && \
		rm $${FILE}
