PROJECT := json-charlieegan3

.PHONY: test
test:
	go test $$(go list ./...)
