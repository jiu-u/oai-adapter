.PHONY: dev
dev:
	go run ./cmd/server/main.go --dev=true

.PHONY: run
run:
	go run./cmd/server/main.go

.PHONY: docker
docker:
	docker build -f ./Dockerfile --build-arg -t zywoo/adapter-demo:v1 .
	docker run --rm -it -p 8080:8080 zywoo/adapter-demo:v1

#docker buildx build --platform linux/amd64 -t zywoo/adapter-demo:v1 .