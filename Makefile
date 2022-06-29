all: test vet fmt lint build

test:
	go test ./...

vet:
	go vet ./...

lint:
	go list ./... | grep -v /vendor/ | xargs -L1 golint -set_exit_status