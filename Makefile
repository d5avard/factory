build-question:
	go build -o question ./cmd/question/.

build-models:
	go build -o models ./cmd/models/.

migration_up:
	migrate -database ${POSTGRESQL_URL} -path database/migration -verbose up

migration_down:
	migrate -database ${POSTGRESQL_URL} -path database/migration -verbose down

migration_fix:
	migrate -database ${POSTGRESQL_URL} -path database/migration/ force VERSION
