.PHONY: build

lint:
	staticcheck ./...

docs:
	godoc -http=:8081

build:
	sam build

tidy:
	go mod tidy && go mod vendor

local:
	sam local start-api

test:
	go test ./... -count=1

coverage:
	go test ./... -coverprofile=cover.out
	go tool cover -html=cover.out

deploy:
	sam deploy