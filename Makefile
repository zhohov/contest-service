include .env
export

# .DEFAULT_GOAL := run

# .PHONY: fmt vet build run

# fmt:
# 	go fmt ./...

# vet: fmt
# 	go vet ./...

# build: vet
# 	go build -o bin/app ./cmd/app

# run: build
# 	./bin/app

# .PHONY: migrate-create
# migrate-create:
#	migrate create -ext sql -dir ./migrations '$(word 2,$(MAKECMDGOALS))'

.PHONY: rebuild
rebuild:
	docker compose up -d --no-deps --build --force-recreate '$(word 2,$(MAKECMDGOALS))'

.PHONY: logs
logs:
	docker compose logs -f '$(word 2,$(MAKECMDGOALS))'

.PHONY: migrate-create
migrate-create:
	docker compose run --rm backend-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq '$(word 2,$(MAKECMDGOALS))'

.PHONY: migrate-up
migrate-up:
	@make migrate-action action=up

.PHONY: migrate-down
migrate-down:
	@make migrate-action action=down

.PHONY: migrate-force
migrate-force:
	@make migrate-action action="force $(v)"

.PHONY: migrate-action
migrate-action:
	docker compose run --rm backend-migrate \
		-path /migrations \
		-database ${DB_URL} \
		'$(action)'
