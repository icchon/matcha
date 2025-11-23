include .env

api-local:
	cd api && go build ./cmd/api-server && ./api-server

fmt:
	gofmt -w .

test-infra:
	@echo "Running DB-dependent tests for store layer..."
	docker-compose down -v
	docker-compose up db -d
	until docker-compose exec db pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}; do \
		sleep 1; \
	done;
	docker-compose run --rm -e TEST_DB_HOST=db api_test_runner sh -c "apk add jq && go mod download && go test -v ./internal/infrastructure/persistence/repo/..." > test-db.log && cat test-db.log || cat test-db.log
	@echo "DB-dependent tests finished."

test-unit:
	@echo "Running DB-independent unit tests..."
	@cd api && \
	go mod tidy && \
	go mod download && \
	PACKAGES=$$(go list ./... | grep -v "/internal/infrastructure" | grep -v "/integration_test"); \
	echo "Packages to test: $$PACKAGES"; \
	if [ -z "$$PACKAGES" ]; then \
		echo "No packages found for unit testing after filtering. This might indicate an issue with the grep patterns or no non-DB/integration tests exist."; \
		exit 1; \
	fi; \
	go test -v $$PACKAGES > ../test-unit.log && cat ../test-unit.log || cat ../test-unit.log
	@echo "DB-independent unit tests finished."

diagram:
	docker-compose up -d db # Ensure db service is running
	until docker-compose exec db pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}; do \
		sleep 1; \
	done;
	docker-compose run --rm tbls doc "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/${POSTGRES_DB}?sslmode=disable" /app/db/tbls --force
