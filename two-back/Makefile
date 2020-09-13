.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/problem problem/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/game/join game/join/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/game/leave game/leave/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/game/publish game/publish/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

