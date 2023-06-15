build:
	@go build -o bin/api
    
run:
	@nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run main.go

seed:
	@go run scripts/seed.go

test:
	@go test -v ./...
