up:
	@goose -dir ./migrations postgres "postgres://avrigne:8679@localhost/explain?sslmode=disable" up
down:
	@goose -dir ./migrations postgres "postgres://avrigne:8679@localhost/explain?sslmode=disable" down
explain:
	@go run explain.go

phony: up down explain