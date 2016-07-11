all: deps format data
		@echo "--> Building Server"
		@go build

deps:
		@echo "--> Installing build dependencies"
		@go get github.com/deckarep/gosx-notifier
		@go get github.com/jteeuwen/go-bindata
		@go get github.com/jteeuwen/go-bindata/go-bindata

format:
		@echo "--> Running go fmt"
		@gofmt -s -w .

data:
		@echo "--> Importing binary files"
		@go-bindata -o scripts.go scripts/