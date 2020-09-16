.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/problem handler/problem/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/join handler/join/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/leave handler/leave/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/solve handler/solve/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

