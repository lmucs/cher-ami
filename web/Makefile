SRC_PATH = src

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
	npm install -g karma-cli