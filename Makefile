.PHONY: build
build:
	docker build . -t faas-autoscaler:latest

.PHONY: deploy
deploy:
	./stack-deploy.sh

.PHONY: clean
clean:
	docker stack rm scaling
