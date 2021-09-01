version = 0.1.0
name = stateservice
local_tag = $(name):$(version)
remote_tag = ralphschaefer/terraform-state-http-service:$(version)

all: build

docker:
	docker build -t $(local_tag) .

docker-publish: docker
	docker tag $(local_tag) $(remote_tag)
	docker push $(remote_tag)

build:
	rm -f $(name)
	CGO_ENABLED=0 GOOS=linux go build -o $(name) -ldflags="-s -w" main.go

run:
	go generate
	go run main.go

clean:
	rm -f main $(name) *~