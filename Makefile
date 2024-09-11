.PHONY: build
build:
	docker compose -f docker-compose.yaml up --build

.PHONY: dev
dev:
	docker compose -f docker-compose.dev.yaml up --build

.PHONY: dev-cleanup
dev-cleanup:
	docker compose -f docker-compose.dev.yaml down
	rm -rf ./sockets/d.sock

.PHONY: test
test:
	sed -i.bak 's|command: \[""\]|command: ["$(PATH_INTEG)"]|' docker-compose.test.yaml
	docker compose -f docker-compose.test.yaml up --build
	mv docker-compose.test.yaml.bak docker-compose.test.yaml
	