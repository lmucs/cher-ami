WEB_SRC = web/src
WEB_ROOT = web
GO_SRC = go/src
GO_TEST_SRC = go/src/test

.PHONY: start
start:
	cd $(GO_SRC); go run main.go

.PHONY: local
local:
	cd $(GO_SRC); go run main.go local

.PHONY: test
test:
	cd $(GO_TEST_SRC); go test -check.v

.PHONY: localtest
localtest:
	cd $(GO_TEST_SRC); go test -check.v -local=true

.PHONY: serve
serve:
	cd $(WEB_SRC); python -m SimpleHTTPServer

.PHONY: watch
watch:
	cd $(WEB_ROOT); compass watch -c config.rb

.PHONY: install-test-reqs
install-test-reqs:
	npm install karma --save-dev
	npm install -g karma-cli

# http://stackoverflow.com/questions/8889035/how-to-document-a-makefile
.PHONY: help
help:
	@echo  'Makefile to assist with initiating the CherAmi project.'
	@echo  ''
	@echo  'Usage:'
	@echo  ''
	@echo  '       make [command]'
	@echo  ''
	@echo  'The commands are:'
	@echo  ''
	@echo  '  start           - Load typical server session using remote database, '
	@echo  '                    and start web project at root path'
	@echo  '  local           - Load server session using local database,'
	@echo  '                    and start web project at root path'
	@echo  '  test            - Start Go server testing using remote database'
	@echo  '  localtest       - Start Go server testing using local database'
	@echo  '  serve           - Serve front-end web application locally to port 8000'
	@echo  '  watch           - Start Compass watcher to keep CSS files up-to-date'
	@echo  ''
