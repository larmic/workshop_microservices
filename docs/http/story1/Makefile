.PHONY: help test

.DEFAULT_GOAL := help

help: ## Zeigt diese Hilfe an
	@echo "Verfügbare Targets:"
	@echo ""
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo ""

test: ## Führt alle HTTP-Tests im docker-Environment aus
	docker run --rm \
	  -v $(CURDIR):/workdir \
	  jetbrains/intellij-http-client \
	  --env using-jetbrains-docker-http-client \
	  --env-file http-client.env.json \
	  requests.http
