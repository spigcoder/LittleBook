.POONY: docker
docker:
	@rm webook || true
	@SET GOOS=linux 
	@SET GOARCH=amd64
	@go build -o webook .
	@docker rmi -f spigcoder/webook:v0.0.1
	@docker build -t spigcoder/webook:v0.0.1 .