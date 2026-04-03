include .env

help:
	echo "MAIN COMMANDS:"
	echo " "
	echo "secrets"
	echo "apply"
	echo "delete"
	echo "---"
	echo "k8s:"
	echo "apply-postgres"
	echo "apply-redis"
	echo "apply-api"
	echo "apply-client"
	echo " "
	echo "delete-postgres"
	echo "delete-redis"
	echo "delete-api"
	echo "delete-client"
	echo "---"
	echo "images:"
	echo "build"
	echo "push"

secrets:
	kubectl apply -f k3s/namespace.yaml
	kubectl apply -f k3s/secrets.yaml

apply-postgres:
	kubectl apply -f k3s/postgres.yaml

delete-postgres:
	kubectl delete -f k3s/postgres.yaml

apply-redis:
	kubectl apply -f k3s/redis.yaml

delete-redis:
	kubectl delete -f k3s/redis.yaml

apply-api:
	kubectl apply -f k3s/api.yaml

delete-api:
	kubectl delete -f k3s/api.yaml

apply-client:
	kubectl apply -f k3s/client.yaml

delete-client:
	kubectl delete -f k3s/client.yaml

apply-ingress:
	kubectl apply -f k3s/https.yaml
	kubectl apply -f k3s/ingress.yaml

delete-ingress:
	kubectl delete -f k3s/https.yaml
	kubectl delete -f k3s/ingress.yaml

apply:
	kubectl apply -f k3s/postgres.yaml
	kubectl apply -f k3s/redis.yaml
	kubectl apply -f k3s/client.yaml
	kubectl apply -f k3s/api.yaml
	kubectl apply -f k3s/ingress.yaml

delete:
	kubectl delete -f k3s/postgres.yaml
	kubectl delete -f k3s/redis.yaml
	kubectl delete -f k3s/api.yaml
	kubectl delete -f k3s/client.yaml
	kubectl delete -f k3s/ingress.yaml

build:
	docker build -t docker-vm-1:5000/books-api:6 ./books-api
	docker build -t docker-vm-1:5000/books-client:4 ./books-client

build-api:
	docker build -t docker-vm-1:5000/books-api:6 ./books-api

build-client:
	docker build -t docker-vm-1:5000/books-client:4 ./books-client

push-client:
	docker push docker-vm-1:5000/books-client:4

push:
	docker push docker-vm-1:5000/books-api:6
	docker push docker-vm-1:5000/books-client:4

deploy:
	docker context use docker-vm-1
	docker compose up -d
