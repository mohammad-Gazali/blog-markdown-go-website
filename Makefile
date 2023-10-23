# you can use here "npm run" or yarn instead of pnpm
build:
	go build -o ./bin/main main.go
	pnpm build

build-dev:
	go build -o ./bin/main main.go

run: build
	./bin/main

run-dev: build-dev
	./bin/main