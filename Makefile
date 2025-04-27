version := $(shell cat ./deployments/web/version.txt)

default: help

# Help target to display available commands
.PHONY: help
help: ## Show this help
	@echo "Usage: make [target]"
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

version: ## Display the current version
	@echo $(version)

migration-up: ## Apply all database migrations
	migrate -database ${POSTGRESQL_URL} -path database/migration -verbose up

migration-down: ## Roll back the last database migration
	migrate -database ${POSTGRESQL_URL} -path database/migration -verbose down

migration-fix: ## Force the migration version
	migrate -database ${POSTGRESQL_URL} -path database/migration/ force VERSION

web-build: ## Build the Docker image for the web service
	docker build -t d5avard/web:$(version) -f ./deployments/web/Dockerfile .

web-push: ## Push the Docker image to the registry
	docker push d5avard/web:$(version)

web-run-debug: ## Run the web service in debug mode using Docker
	docker run -d \
		--user webuser \
		-p 8080:80 \
		-p 8443:443 \
		-e WEB_CERT_FILE="/app/certs/fullchain.pem" \
		-e WEB_KEY_FILE="/app/certs/privkey.pem" \
		-v /Users/danysavard/Projects/factory/web/web/certs:/app/certs:ro \
		d5avard/web:$(version)

web-run-prod: ## Run the web service in production mode using Docker
	docker run -d \
		--user webuser \
		-p 80:80 \
		-e WEB_CERT_FILE="/app/certs/fullchain.pem" \
		-e WEB_KEY_FILE="/app/certs/privkey.pem" \
		-v /Users/danysavard/Projects/factory/web/web/certs:/app/certs:ro \
		-p 443:443 \
		d5avard/web:$(version)

web-run-local: ## Run the web service locally for development
	WEB_CERT_FILE=./web/web/certs/fullchain.pem \
	WEB_KEY_FILE=./web/web/certs/privkey.pem \
		go run ./cmd/web/web.go \
		--host localhost \
		--httpPort 8080 \
		--tlsPort 8443 \
		--templatePath "./web/web/templates"


