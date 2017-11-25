run:
	go build -o build . && ./build

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

host:=http://localhost:8080
auth:=admin
token:=token

curl-auth-login:
	curl $(host)/api/auth/login?token=$(token)

curl-members-id:
	curl -H 'Auth:$(auth)' $(host)/api/members/$(id)

curl-members:
	curl -H 'Auth:$(auth)' $(host)/api/members
