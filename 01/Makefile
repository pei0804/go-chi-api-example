run:
	go build -o build . && ./build

fresh:
	fresh

host:=http://localhost:8080
auth:=admin
token:=token

curl-auth-login:
	curl $(host)/api/auth/login?token=$(token)

curl-members-id:
	curl -H 'Authorization:$(auth)' $(host)/api/members/$(id)

curl-members:
	curl -H 'Authorization:$(auth)' $(host)/api/members
