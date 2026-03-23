.PHONY: help all build clean run-flight run-hotel run-car run-booking-story1 run-all test docker-build docker-build-flight docker-build-hotel docker-build-car docker-build-booking-story1 docker-up docker-up-hub docker-down

.DEFAULT_GOAL := help

help: ## Zeigt diese Hilfe an
	@echo "Verfügbare Targets:"
	@echo ""
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo ""

all: build ## Baut alle Services

build: ## Baut alle Service-Binaries nach bin/
	go build -o bin/flight ./flight
	go build -o bin/hotel ./hotel
	go build -o bin/car ./car
	go build -o bin/booking-story1 ./booking/story1

clean: ## Löscht das bin/ Verzeichnis
	rm -rf bin/

run-flight: ## Startet den Flight-Service lokal
	go run ./flight

run-hotel: ## Startet den Hotel-Service lokal
	go run ./hotel

run-car: ## Startet den Car-Service lokal
	go run ./car

run-booking-story1: ## Startet den BookingService Story1 lokal
	go run ./booking/story1

run-all: ## Startet alle Services lokal (parallel)
	@echo "Starting all services..."
	@go run ./flight & go run ./hotel & go run ./car & go run ./booking/story1 & wait

test: ## Testet alle laufenden Services via curl
	@echo "Testing FlightService..."
	@curl -s http://localhost:8081/health && echo ""
	@curl -s http://localhost:8081/openapi | head -5
	@echo "\nTesting HotelService..."
	@curl -s http://localhost:8082/health && echo ""
	@curl -s http://localhost:8082/openapi | head -5
	@echo "\nTesting CarService..."
	@curl -s http://localhost:8083/health && echo ""
	@curl -s http://localhost:8083/openapi | head -5
	@echo "\nTesting BookingService Story1..."
	@curl -s http://localhost:8085/health && echo ""
	@curl -s http://localhost:8085/openapi | head -5

docker-build: docker-build-flight docker-build-hotel docker-build-car docker-build-booking-story1 ## Baut alle Docker-Images

docker-build-flight: ## Baut das Flight-Service Docker-Image
	docker build -f flight/Dockerfile -t workshop-microservices-flight:latest .

docker-build-hotel: ## Baut das Hotel-Service Docker-Image
	docker build -f hotel/Dockerfile -t workshop-microservices-hotel:latest .

docker-build-car: ## Baut das Car-Service Docker-Image
	docker build -f car/Dockerfile -t workshop-microservices-car:latest .

docker-build-booking-story1: ## Baut das BookingService Story1 Docker-Image
	docker build -f booking/story1/Dockerfile -t workshop-microservices-booking-story1:latest .

docker-up: ## Baut und startet alle Services mit Docker Compose
	docker compose -f docker-compose.yml -f docker-compose.infra.yml up --build

docker-up-hub: ## Zieht Images von Docker Hub und startet
	docker compose -f docker-compose.yml -f docker-compose.infra.yml pull
	docker compose -f docker-compose.yml -f docker-compose.infra.yml up

docker-down: ## Stoppt und entfernt alle Container
	docker compose down
