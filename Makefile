.PHONY: help test up clean swagger

help:
	@echo "Available targets:"
	@echo "  help        - Show this help"
	@echo "  test        - Run all tests with race detector and coverage"
	@echo "  up          - Start Docker containers"
	@echo "  clean       - Clean build artifacts"
	@echo "  swagger     - Generate Swagger documentation"

test:
	go get github.com/stretchr/testify/assert
	go get github.com/stretchr/testify/mock
	go mod tidy
	go test -v -race -coverprofile=coverage.out ./...

up:
	docker-compose up -d

clean:
	rm -rf bin/
	rm -f coverage.out
	rm -rf internal/api/rest/swagger/docs/

swagger:
	@echo "Installing swag if not present..."
	@which swag > /dev/null || go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Generating Swagger documentation..."
	cd internal/api/rest/handlers && swag init --parseDependency --generalInfo ../server.go --output ../swagger/docs/
	@echo "Swagger documentation generated successfully!"
	@echo "Access it at http://localhost:8080/swagger/index.html"
