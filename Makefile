include .env

fmt:
	gofmt -w .

up:
	docker-compose up -d

down:
	docker-compose down

downv:
	docker-compose down -v

tool:
	docker-compose --profile tools up -d

diagram:
	docker-compose up -d db # Ensure db service is running
	until docker-compose exec db pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}; do \
		sleep 1; \
	done;
	docker-compose run --rm tbls doc "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/${POSTGRES_DB}?sslmode=disable" /app/db/tbls --force
	docker-compose down

seed:
	docker-compose up -d db # Ensure db service is running
	until docker-compose exec db pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}; do \
		sleep 1; \
	done;
	docker-compose build seeder
	docker-compose run --rm seeder
	docker-compose down

web-install:
	cd web && npm install

web-dev:
	cd web && npm run dev

web-build:
	cd web && npm run build

web-test:
	cd web && npm run test:run

web-lint:
	cd web && npm run lint

web-format:
	cd web && npm run format
