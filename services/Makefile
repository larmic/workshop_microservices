.PHONY: all build clean run-flight run-hotel run-car run-all test

all: build

build:
	go build -o bin/flight ./flight
	go build -o bin/hotel ./hotel
	go build -o bin/car ./car

clean:
	rm -rf bin/

run-flight:
	go run ./flight

run-hotel:
	go run ./hotel

run-car:
	go run ./car

run-all:
	@echo "Starting all services..."
	@go run ./flight & go run ./hotel & go run ./car & wait

test:
	@echo "Testing FlightService..."
	@curl -s http://localhost:8081/health && echo ""
	@curl -s http://localhost:8081/openapi | head -5
	@echo "\nTesting HotelService..."
	@curl -s http://localhost:8082/health && echo ""
	@curl -s http://localhost:8082/openapi | head -5
	@echo "\nTesting CarService..."
	@curl -s http://localhost:8083/health && echo ""
	@curl -s http://localhost:8083/openapi | head -5
