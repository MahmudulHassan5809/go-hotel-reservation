build:
	@go build -o bin/api
    
run:
	@nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run main.go

seed:
	@go run scripts/seed.go

docker:
	echo "building docker file"
	@docker build -t api .
	echo "running API inside Docker container"
	@docker run -p 8000:8000 api

test:
	@go test -v ./...
