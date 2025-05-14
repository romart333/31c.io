# docker commands

docker-up:
	docker start db-test || docker run -d -p 3406:3306 --env MYSQL_ROOT_PASSWORD=root_password --name db-test --rm mysql:latest

docker-stop:
	docker stop db-test

# app commands

start:
	go run .

start-dev:
	reflex -r '\.go$$' -R '^vendor/' -R '^docs/' -s -- bash -c "swag init && go run ."

go-tidy:
	go mod tidy

go-vendor:
	 go mod vendor

swag-init:
	swag init

test-run:
	go test

test-run-v:
	go test -v

test-run-s:
	go test -run TestAllRoutes/TestAccountRoute/TestGetAccountsRoute_SuccessNoParams
