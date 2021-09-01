version = 0.1.0
name = stateservice
local_tag = $(name):$(version)

all: build

docker:
	docker build -t $(local_tag) .

build:
	rm -f $(name)
	CGO_ENABLED=0 GOOS=linux go build -o $(name) -ldflags="-s -w" main.go

run:
	go generate
	go run main.go

clean:
	rm -f main $(name) *~