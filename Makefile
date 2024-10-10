run:
	go run main.go

migrate_up:
	migrate -database postgres://postgres:postgres@localhost:5432/go_practice?sslmode=disable -path db/migrations up

migrate_down:
	migrate -database postgres://postgres:postgres@localhost:5432/go_practice?sslmode=disable -path db/migrations down

migrate:
	migrate create -ext sql -dir db/migrations -seq db

das