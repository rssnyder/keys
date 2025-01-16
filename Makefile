build:
	docker build -t registry.rileysnyder.dev/keys:test .

test:
	PORT=3948 docker compose up --build -d
	go test
	docker compose rm -fsv

push:
	docker push registry.rileysnyder.dev/keys:test