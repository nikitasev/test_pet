COMPOSE_INFRA_FILE_DEV=deployments/local/infrastructure.yaml
COMPOSE_APP_FILE_DEV=deployments/local/app.yaml
COMPOSE_PROJECT_APP=pet_service
COMPOSE_PROJECT_INFRA=pet_infra
DOCKER_NETWORK_NAME=pet_network

run-infra:
	@echo "export COMPOSE_FILE=$(COMPOSE_INFRA_FILE_DEV); export COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT)"
	docker network ls|grep $(DOCKER_NETWORK_NAME) > /dev/null || docker network create $(DOCKER_NETWORK_NAME) && \
	docker-compose -f "$(COMPOSE_INFRA_FILE_DEV)" -p "$(COMPOSE_PROJECT_INFRA)" down --remove-orphans && \
	docker-compose -f "$(COMPOSE_INFRA_FILE_DEV)" -p "$(COMPOSE_PROJECT_INFRA)" up -d --build --force-recreate
