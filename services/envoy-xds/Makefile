.PHONY: build



#@env GOOS=linux GOARCH=amd64 go build

build:
	@echo "--> building for Linux platform"
	@env GOOS=linux GOARCH=amd64 go build  -mod vendor

.PHONY: clean
clean:
	@echo "-->clean"
	@rm -rf envoy-xds
