run:
	go build . && ./goapi

test:
	env HTTPDOC=1 go test -v ./...

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-stop:
	docker-compose stop

docker-rm:
	docker-compose rm

docker-ssh:
	docker exec -it goapi /bin/bash

docker-server: docker-build docker-up

docker-clean: docker-stop docker-rm
