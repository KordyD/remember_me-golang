all: run

run:
	go run cmd/sso/main.go --config=config/config.yaml

test: test-unit

test-unit:
	go test $$(go list ./... | grep -v /tests/ | grep -v /mocks) -count=1 -v

test-integration:
	go test $$(go list ./tests/integration/...) -count=1 -v

