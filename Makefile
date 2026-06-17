up:
	docker compose up --build -d
down:
	docker compose down
logs:
	docker compose logs -f user-service
ps:
	docker compose ps
migrate:
	docker compose exec postgres psql -U postgres -d userdb -f /migrations/000001_create_users_table.up.sql