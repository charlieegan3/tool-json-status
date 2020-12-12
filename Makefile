.PHONY: test

PROJECT := json-charlieegan3
TAG := $(shell tar -cf - . | md5sum | cut -f 1 -d " ")

test:
	go test $$(go list ./...)

build:
	docker build -t charlieegan3/$(PROJECT):latest -t charlieegan3/$(PROJECT):${TAG} .

push: build
	docker push charlieegan3/$(PROJECT):latest
	docker push charlieegan3/$(PROJECT):${TAG}

build_arm:
	docker build -t charlieegan3/$(PROJECT):arm-${TAG} -f Dockerfile.arm .

push_arm: build_arm
	docker push charlieegan3/$(PROJECT):arm-${TAG}
