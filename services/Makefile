.PHONY: help all build clean run-flight run-hotel run-car run-booking-story1 run-booking-story2 run-all test docker-build docker-build-flight docker-build-hotel docker-build-car docker-build-booking-story1 docker-build-booking-story2 docker-up docker-up-base docker-up-hub docker-down

.DEFAULT_GOAL := help

help: ## Zeigt diese Hilfe an
	@echo "Verfügbare Targets:"
	@echo ""
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo ""

docker-build: docker-build-flight docker-build-hotel docker-build-car docker-build-booking-story1 docker-build-booking-story2 ## Baut alle Docker-Images

docker-build-flight: ## Baut das Flight-Service Docker-Image
	docker build -f flight/Dockerfile -t workshop-microservices-flight:latest .

docker-build-hotel: ## Baut das Hotel-Service Docker-Image
	docker build -f hotel/Dockerfile -t workshop-microservices-hotel:latest .

docker-build-car: ## Baut das Car-Service Docker-Image
	docker build -f car/Dockerfile -t workshop-microservices-car:latest .

docker-build-booking-story1: ## Baut das BookingService Story1 Docker-Image
	docker build -f booking/story1/Dockerfile -t workshop-microservices-booking-story1:latest .

docker-build-booking-story2: ## Baut das BookingService Story2 Docker-Image
	docker build -f booking/story2/Dockerfile -t workshop-microservices-booking-story2:latest .

docker-up: ## Baut und startet alle Services mit Docker Compose
	docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml up --build

docker-up-hub: ## Zieht Images von Docker Hub und startet alle Services
	docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml pull
	docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml up
