.POONY: docker
docker:
	@rm webook || true
	@SET GOOS=linux 
	@SET GOARCH=amd64
	@go build -o webook .
	@docker rmi -f spigcoder/webook:v0.0.1
	@docker build -t spigcoder/webook:v0.0.1 .

.PHONY: mock
mock:
	@mockgen -source=internal/repository/article/article_author.go -destination=internal/repository/article/mocks/article_author_mock.go -package=artRepoMocks
	@mockgen -source=internal/repository/article/article_reader.go -destination=internal/repository/article/mocks/article_reader_mock.go -package=artRepoMocks
	@go mod tidy