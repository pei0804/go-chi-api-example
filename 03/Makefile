run:
	go build -o build . && ./build

test:
	export HTTPDOC=1 && go test -v ./...

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

config-set: 
	cp -i db.yml.tmpl db.yml

host:=http://localhost:8080
auth:=admin
token:=token

curl-auth-login:
	curl $(host)/api/auth/login?token=$(token)

curl-members-id:
	curl -H 'Authorization:$(auth)' $(host)/api/members/$(id)

curl-members:
	curl -H 'Authorization:$(auth)' $(host)/api/members

curl-reqid:
	curl -H 'Authorization:$(auth)' -H 'X-Request-ID:test' $(host)/api/members
