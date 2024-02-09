.PHONY: run-tests
run-tests:
	@go test -v -failfast `go list ./...` -cover

.PHONY:
mock-install:
	@go install go.uber.org/mock/mockgen@v0.4.0

.PHONY: mock
mock:
	@`go env GOPATH`/bin/mockgen -source ./$(util)/$(subutil).go -destination ./tests/mock/$(util)/$(subutil).go

.PHONY: mock-all
mock-all:
	@make mock util=auth subutil=auth
	@make mock util=configbuilder subutil=configbuilder
	@make mock util=configreader subutil=configreader
	@make mock util=email subutil=email
	@make mock util=email subutil=email_template
	@make mock util=featureflag subutil=feature_flag
	@make mock util=log subutil=log
	@make mock util=messaging subutil=messaging
	@make mock util=parser subutil=parser
	@make mock util=parser subutil=csv
	@make mock util=parser subutil=json
	@make mock util=ratelimiter subutil=rate_limiter
	@make mock util=redis subutil=redis
	@make mock util=security subutil=security
