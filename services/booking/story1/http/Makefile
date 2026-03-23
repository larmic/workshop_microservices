.PHONY: help test

.DEFAULT_GOAL := help

UNAME_S := $(shell uname -s)
DOCKER_HOST_FLAG :=
ifeq ($(UNAME_S),Linux)
  DOCKER_HOST_FLAG := --add-host host.docker.internal:host-gateway
endif

help: ## Zeigt diese Hilfe an
	@echo "Verfügbare Targets:"
	@echo ""
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo ""

test: ## Führt alle HTTP-Tests im docker-Environment aus
	docker run --rm \
	  $(DOCKER_HOST_FLAG) \
	  -v $(CURDIR):/workdir \
	  jetbrains/intellij-http-client \
	  --env docker \
	  --env-file http-client.env.json \
	  requests.http
