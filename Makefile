WEB_SRC = web/src
WEB_ROOT = web
GO_SRC = go/src

.PHONY: start
start:
	cd $(GO_SRC); go run main.go

.PHONY: serve
serve:
	cd $(WEB_SRC); python -m SimpleHTTPServer

.PHONY: watch
watch:
	cd $(WEB_ROOT); compass watch -c config.rb

install-test-reqs:
	npm install karma --save-dev
	npm install -g karma-cli
