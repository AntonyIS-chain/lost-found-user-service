build:
	go build -o bin/backend/user-service
	
serve: build
	ENV=development ./bin/backend/user-service

serve-dev-test: build
	ENV=development_test go test -v ./...

docker-push:
	docker build -t antonyinjila/backend/user-service:latest --build-arg ENV=docker .
	docker push antonyinjila/backend/user-service:latest

docker-run:
	docker run -p 8081:8081 ENV=docker antonyinjila/backend/user-service:latest

docker-test:
	ENV=docker_test go test -v ./...