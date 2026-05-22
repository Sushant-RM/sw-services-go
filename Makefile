APP_NAME=sw-services

up:
	docker-compose up --build

down:
	docker-compose down

clean:
	docker-compose down -v

build:
	go build ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

check:
	go fmt ./...
	go vet ./...
	go test ./...
	go build ./...

logs:
	docker-compose logs -f

restart:
	docker-compose down
	docker-compose up --build

health:
	curl http://localhost:3468/sw-services/actuator/health
