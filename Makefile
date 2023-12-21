build:
	@echo "Building and run the"
	go build

run: build
	@echo "Run base-service"
	./base-service