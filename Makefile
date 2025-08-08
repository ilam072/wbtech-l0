up:
	docker-compose up -d
	go run ./frontend/main.go &
	go run ./backend/cmd/server/main.go
down:
	docker-compose down

.PHONY: producer
producer:
	docker compose exec kafka kafka-console-producer.sh --bootstrap-server localhost:9092 --topic $(topic)

topics:
	docker exec -it kafka kafka-topics.sh --bootstrap-server localhost:9092 --list

messages:
	docker exec -it kafka kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic orders --from-beginning