dev:
	@nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run main.go

build:
	@go build -o secretariat_repository

run:
	@./secretariat_repository

build_run: build run