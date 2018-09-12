.PHONY: build


#@env GOOS=linux GOARCH=amd64 go build  -mod vendor

build:
	@echo "--> building for Linux platform"
	@env GOOS=linux GOARCH=amd64 go build

.PHONY: clean
clean:
	@echo "-->clean"
	@rm -rf envoy-xds
