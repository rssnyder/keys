test:
	docker compose rm -fsv
	PORT=3948 docker compose up --build -d
	go test