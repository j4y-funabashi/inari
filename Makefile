test: test-ui test-api test-acceptance

test-api:
	docker-compose run --rm api ./test.sh

test-ui:
	docker-compose run --rm ui ./test.sh

test-acceptance:
	docker-compose run --rm acceptance firefox:headless
