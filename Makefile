build:
	docker build -t docker-vm-1:5000/books-api:latest ./books-api
	docker build -t docker-vm-1:5000/books-client:latest ./books-client

push:
	docker push docker-vm-1:5000/books-api:latest
	docker push docker-vm-1:5000/books-client:latest

deploy:
	docker context use docker-vm-1
	docker compose up -d
