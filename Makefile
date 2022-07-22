test: test-ui test-api

test-api:
	docker compose run --rm api ./test.sh

test-ui:
	docker compose run --rm ui ./test.sh

test-acceptance:
	docker compose run --rm acceptance firefox:headless

deploy:
	docker-compose run --rm deployer ./apps/deployer/deploy.sh
