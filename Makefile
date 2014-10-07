WEB_SRC = web/src
WEB_ROOT = web
GO_SRC = go/src

.PHONY: start
start:
	cd $(GO_SRC); go run main.go "http://107.170.229.205:7474/db/data"

.PHONY: local
local:
	cd $(GO_SRC); go run main.go "http://localhost:7474/db/data"

.PHONY: test-api
test-api:
	cd $(GO_SRC); go run main.go "http://192.241.226.228:7474/db/data"

.PHONY: serve
serve:
	cd $(WEB_SRC); python -m SimpleHTTPServer

.PHONY: watch
watch:
	cd $(WEB_ROOT); compass watch -c config.rb

install-test-reqs:
	npm install karma --save-dev
	npm install -g karma-cli
