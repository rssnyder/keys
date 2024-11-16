test:
	docker compose rm -fsv
	docker compose up --build -d
	go test