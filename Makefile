WEB_ROOT        = web
WEB_SRC         = web/src
API_SERVER      = api-server
API_TEST_DIR    = api-server/test

.PHONY: start
start:
	cd $(API_SERVER); go run api-server.go

.PHONY: local
local:
	cd $(API_SERVER); go run api-server.go local

.PHONY: test
test:
	cd $(API_TEST_DIR); go test -check.v

.PHONY: localtest
localtest:
	cd $(API_TEST_DIR); go test -check.v -local=true

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

.PHONY: deps
install-deps:
	@echo '--------------------------'
	@echo 'Getting Go dependencies...'
	@echo '--------------------------'
	cd $(API_SERVER); go get all; cd ..
	@echo '-----------------'
	@echo 'Packages installed.'
	@echo '-----------------'

.PHONY: update-deps
update-deps:
	@echo '---------------------------'
	@echo 'Updating Go dependencies...'
	@echo '---------------------------'
	cd $(API_SERVER); go get -u all; cd ..
	@echo '-----------------'
	@echo 'Packages updated.'
	@echo '-----------------'

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
	@echo  '  install-deps    - Install Go package dependencies'
	@echo  '  update-deps     - Update Go package dependencies'
	@echo  ''
