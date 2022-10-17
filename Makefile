IMAGE_REPO ?= ""
IMAGE_NAME ?= injector-secret
IMAGE_TAG ?= latest

image: build-image push-image

build-image:
	@docker build -t \
		$(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG) \
		-f build/Dockerfile ./cmd

push-image:
	@docker push $(IMAGE_REPO)/$(IMAGE_NAME):latest

install:
	helm upgrade -i $(IMAGE_NAME) ./chart

delete:
	kubectl delete secret $(IMAGE_NAME)
	helm delete $(IMAGE_NAME)

.PHONY: all image install delete
