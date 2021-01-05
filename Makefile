# test dependencies for known vulnerable paths
test_security:
	snyk test

test: test_security
	go test ./... -short -coverprofile=coverage.out
