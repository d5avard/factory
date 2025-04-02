version := $(shell cat ./deployments/web/version.txt)

build-question:
	go build -o question ./cmd/question/.

build-models:
	go build -o models ./cmd/models/.

migration-up:
	migrate -database ${POSTGRESQL_URL} -path database/migration -verbose up

migration-down:
	migrate -database ${POSTGRESQL_URL} -path database/migration -verbose down

migration-fix:
	migrate -database ${POSTGRESQL_URL} -path database/migration/ force VERSION

web-build:
	docker build -t d5avard/web:$(version) -f ./deployments/web/Dockerfile .

web-push:
	docker push d5avard/web:$(version)

web-run-local:
	go run ./cmd/web/web.go --host localhost --port 8080 --templatePath "./web/web/templates"

web-run-debug:
	docker run -d --user webuser -p 8080:80 d5avard/$(version)