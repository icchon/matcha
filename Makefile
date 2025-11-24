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
