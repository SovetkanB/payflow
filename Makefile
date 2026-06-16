up:
	docker compose up --build -d
down:
	docker compose down
logs:
	docker compose logs -f user-service
ps:
	docker compose ps