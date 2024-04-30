app.run:
	docker-compose up

test:
	docker run -it --rm -v .:/app -v /var/run/docker.sock:/var/run/docker.sock testcontainers-example-app go test ./... -v