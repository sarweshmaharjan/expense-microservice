DOCKER_COMPOSE_FILE=~/code/work/sarweshmaharjan/expenses/hub/docker_compose.yml

.PHONY: build start stop remove rebuild

build:
	docker-compose -f $(DOCKER_COMPOSE_FILE) build

start:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

stop:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

remove:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down --rmi all --volumes --remove-orphans

rebuild:
	$(MAKE) stop
	$(MAKE) remove
	$(MAKE) start