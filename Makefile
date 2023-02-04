vet:
	go vet ./...

fmt:
	gofmt -s -w ./
	goimports -w ./