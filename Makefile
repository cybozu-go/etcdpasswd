all:
	go install ./...

test: setup_test
	go test -race -v ./...

setup_test: mocks/auth.go

mocks/auth.go: auth.go
	mkdir -p mocks
	go get github.com/golang/mock/gomock
	go get github.com/golang/mock/mockgen
	go generate ./...
	goimports -w $@

.PHONY: all test setup_test
