.PHONY: all build clean run-flight run-hotel run-car run-all test docker-build docker-build-flight docker-build-hotel docker-build-car

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

docker-build: docker-build-flight docker-build-hotel docker-build-car

docker-build-flight:
	docker build -f flight/Dockerfile -t workshop-microservices-flight:latest .

docker-build-hotel:
	docker build -f hotel/Dockerfile -t workshop-microservices-hotel:latest .

docker-build-car:
	docker build -f car/Dockerfile -t workshop-microservices-car:latest .
