.PHONY: build
build:
	docker compose -f docker-compose.dev.yaml up --build

.PHONY: test
test:
	sed -i.bak 's|command: \[""\]|command: ["$(PATH_INTEG)"]|' docker-compose.test.yaml
	docker compose -f docker-compose.test.yaml up --build
	mv docker-compose.test.yaml.bak docker-compose.test.yaml