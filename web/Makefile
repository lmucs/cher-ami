SRC_PATH = static

.PHONY: all
all: serve watch

.PHONY: serve
serve:
	cd $(SRC_PATH); python -m SimpleHTTPServer

.PHONY: watch
watch:
	compass watch -c config.rb

install-test-reqs:
	npm install karma --save-dev
	npm install karma-jasmine@2_0 karma-chrome-launcher --save-dev
	npm install karma-requirejs --save-dev

install-test-cli:
	npm install -g karma-cli