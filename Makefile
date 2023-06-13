build:
	@go build -o bin/api
    
run:
	@nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run main.go
	

test:
	@go test -v ./...
