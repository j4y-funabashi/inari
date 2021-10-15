.PHONY: build clean deploy gomodgen

build: gomodgen
	./build_lambda.sh

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	./node_modules/serverless/bin/serverless.js deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
