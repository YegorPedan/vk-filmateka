ABS_PATH=$(CURDIR)

test:
	CONFIG_PATH=$(ABS_PATH)/config/test.yaml go test ./...

cover:
	CONFIG_PATH=$(ABS_PATH)/config/test.yaml go test ./... -v -coverpkg=./... -coverprofile=c.out
	go tool cover -html="c.out"
	rm c.out

swagger:
	rm -rf ./docs
	swag init -g cmd/main/main.go

run:
	CONFIG_PATH=./config/local.yaml go run ./cmd/main/main.go

dev:
	CONFIG_PATH=./config/dev.yaml go run ./cmd/main/main.go

prod:
	CONFIG_PATH=./config/prod.yaml go run ./cmd/main/main.go

clean:
	rm -rf main

build:
	CONFIG_PATH=./config/prod.yaml go build ./cmd/main/main.go

start:
	CONFIG_PATH=./config/prod.yaml ./main