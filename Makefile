.PHONY: dev
dev:
	go run ./cmd/server/ --dev=true

.PHONY: run
run:
	go run./cmd/server/

.PHONY: docker
docker:
	docker build -f ./Dockerfile --build-arg -t zywoo/adapter-demo:v1 .
	docker run --rm -it -p 8080:8080 zywoo/adapter-demo:v1

#docker buildx build --platform linux/arm64 -t zywoo/adapter-demo-arm64:v1 .
# docker save zywoo/adapter-demo-arm64:v1 -o dist.tar