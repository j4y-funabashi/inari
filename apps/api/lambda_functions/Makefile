.PHONY: build clean deploy gomodgen

build:
	./build_lambda.sh

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	./node_modules/serverless/bin/serverless.js deploy --verbose
