.PHONY: help all build clean run-flight run-hotel run-car run-booking-story1 run-booking-story2 run-booking-story3 run-booking-story4 run-booking-story5 run-booking-story6 run-booking-story7 run-all test docker-build docker-build-flight docker-build-hotel docker-build-car docker-build-booking-story1 docker-build-booking-story2 docker-build-booking-story3 docker-build-booking-story4 docker-build-booking-story5 docker-build-booking-story6 docker-build-booking-story7 docker-build-dashboard docker-up docker-up-base docker-up-hub docker-down

.DEFAULT_GOAL := help

help: ## Zeigt diese Hilfe an
	@echo "Verfügbare Targets:"
	@echo ""
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo ""

docker-build: docker-build-flight docker-build-hotel docker-build-car docker-build-booking-story1 docker-build-booking-story2 docker-build-booking-story3 docker-build-booking-story4 docker-build-booking-story5 docker-build-booking-story6 docker-build-booking-story7 docker-build-dashboard ## Baut alle Docker-Images

docker-build-flight: ## Baut das Flight-Service Docker-Image
	docker build -f Dockerfile --build-arg SERVICE_PATH=flight --build-arg SERVICE_DESC="Workshop Microservices - Flight Service" --build-arg SERVICE_PORT=8081 -t workshop-microservices-flight:latest .

docker-build-hotel: ## Baut das Hotel-Service Docker-Image
	docker build -f Dockerfile --build-arg SERVICE_PATH=hotel --build-arg SERVICE_DESC="Workshop Microservices - Hotel Service" --build-arg SERVICE_PORT=8082 -t workshop-microservices-hotel:latest .

docker-build-car: ## Baut das Car-Service Docker-Image
	docker build -f Dockerfile --build-arg SERVICE_PATH=car --build-arg SERVICE_DESC="Workshop Microservices - Car Service" --build-arg SERVICE_PORT=8083 -t workshop-microservices-car:latest .

docker-build-booking-story1: ## Baut das BookingService Story1 Docker-Image
	docker build -f Dockerfile --build-arg SERVICE_PATH=booking/story1 --build-arg SERVICE_DESC="Workshop Microservices - Booking Service (Story 1)" --build-arg SERVICE_PORT=8080 -t workshop-microservices-booking:story1 .

docker-build-booking-story2: ## Baut das BookingService Story2 Docker-Image
	docker build -f Dockerfile --build-arg SERVICE_PATH=booking/story2 --build-arg SERVICE_DESC="Workshop Microservices - Booking Service (Story 2)" --build-arg SERVICE_PORT=8080 -t workshop-microservices-booking:story2 .

docker-build-booking-story3: ## Baut das BookingService Story3 Docker-Image
	docker build -f Dockerfile --build-arg SERVICE_PATH=booking/story3 --build-arg SERVICE_DESC="Workshop Microservices - Booking Service (Story 3)" --build-arg SERVICE_PORT=8080 -t workshop-microservices-booking:story3 .

docker-build-booking-story4: ## Baut das BookingService Story4 Docker-Image
	docker build -f Dockerfile --build-arg SERVICE_PATH=booking/story4 --build-arg SERVICE_DESC="Workshop Microservices - Booking Service (Story 4)" --build-arg SERVICE_PORT=8080 -t workshop-microservices-booking:story4 .

docker-build-booking-story5: ## Baut das BookingService Story5 Docker-Image
	docker build -f Dockerfile --build-arg SERVICE_PATH=booking/story5 --build-arg SERVICE_DESC="Workshop Microservices - Booking Service (Story 5)" --build-arg SERVICE_PORT=8080 -t workshop-microservices-booking:story5 .

docker-build-booking-story6: ## Baut das BookingService Story6 Docker-Image
	docker build -f Dockerfile --build-arg SERVICE_PATH=booking/story6 --build-arg SERVICE_DESC="Workshop Microservices - Booking Service (Story 6)" --build-arg SERVICE_PORT=8080 -t workshop-microservices-booking:story6 .

docker-build-booking-story7: ## Baut das BookingService Story7 Docker-Image
	docker build -f Dockerfile --build-arg SERVICE_PATH=booking/story7 --build-arg SERVICE_DESC="Workshop Microservices - Booking Service (Story 7)" --build-arg SERVICE_PORT=8080 -t workshop-microservices-booking:story7 .

docker-build-dashboard: ## Baut das Scaling-UI Docker-Image
	docker build -f dashboard/Dockerfile -t workshop-microservices-dashboard:latest .

docker-up: ## Baut und startet alle Services mit Docker Compose
	docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml up --build

docker-up-hub: ## Zieht Images von Docker Hub und startet alle Services
	docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml pull
	docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml up
