test:
	@go test -v -race -coverprofile cov.out -cpuprofile cpu.out -memprofile mem.out
	@go tool cover -html cov.out
