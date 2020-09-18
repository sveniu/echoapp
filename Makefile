.docker-build:
	docker build -t echoapp .

.docker-tag: .docker-build
	docker tag echoapp:latest 465509910691.dkr.ecr.eu-central-1.amazonaws.com/echoapp:latest

.docker-push: .docker-tag
	docker push 465509910691.dkr.ecr.eu-central-1.amazonaws.com/echoapp:latest

all: .docker-push
