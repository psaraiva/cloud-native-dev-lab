GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
SWAG=swag
MOCKERY=mockery

all: build

.PHONY: run
run:
	@echo "=> Runing temporary app..."
	$(GOCMD) run .

.PHONY: run-build
run-build: clean build
	@echo "=> Run by build..."
	./server

.PHONY: build
build:
	@echo "=> Build the app..."
	$(GOBUILD) -o server .

.PHONY: cyclo
cyclo:
	@echo "Cyclomatic Complexity: Top 10..."
	gocyclo -top 10 ./service ./handler ./client

.PHONY: lint
lint:
	@echo "Analyzing..."
	golangci-lint run --tests=true ./service ./client

.PHONY: mock
mock:
	@echo "Generating mocks"
	rm -rf mocks/
	$(MOCKERY) --name=PokemonProvider --dir=./service --output=./mocks --unroll-variadic=true

.PHONY: test
test:
	@echo "Runing unit tests..."
	$(GOTEST) -v -cover ./...

.PHONY: quality-metrics
quality-metrics: cyclo lint

.PHONY: docs
docs:
	@echo "Generating documentation by swag..."
	rm -rf docs/
	$(SWAG) init

.PHONY: clean
clean:
	@echo Broom": SHHHHH... SHHHH..."
	rm -f server coverage.out
	rm -rf mocks/
