
.PHONY: run
run:
	@echo "Running..."
	@swag init
	@docker-compose up --build

.PHONY: install
install:
	@echo "Installing..."
	@go install
	@go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: clean
clean:
	@echo "Cleaning..."
	@docker-compose down -v
	@rm -rf docs

.PHONY: test
test:
	@echo "Testing..."
	@go test -v ./...

.PHONY: bench
bench:
	@echo "Benchmarking..."
	@go test -bench=. ./...
