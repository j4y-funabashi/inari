build: build-api build-ui build-cli

build-cli:
	docker build -t inari-cli --file ./apps/api/Dockerfile-cli ./apps/api

build-ui:
	docker compose build inari-ui

build-api:
	docker compose build inari-api

up:
	docker compose up --remove-orphans --detach

down:
	docker compose down

test: test-ui test-api

test-api:
	docker compose run --rm inari-api-test ./test.sh

test-ui:
	docker compose run --build inari-ui-test npm test

test-acceptance:
	docker compose run --rm acceptance firefox:headless

deploy:
	docker-compose run --rm deployer ./apps/deployer/deploy.sh
