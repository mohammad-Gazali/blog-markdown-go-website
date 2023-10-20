# you can use here "npm run" or yarn instead of pnpm
build:
	go build -o ./bin/main main.go
	pnpm build

run: build
	./bin/main