up:
	docker-compose up -d
#	go run ./frontend/main.go &
#	go run ./backend/cmd/app/main.go
down:
	docker-compose down

.PHONY: producer
producer:
	docker compose exec kafka-l0 kafka-console-producer.sh --bootstrap-server kafka-l0:9092 --topic $(topic)

topics:
	docker exec -it kafka-l0 kafka-topics.sh --bootstrap-server localhost:9092 --list

messages:
	docker exec -it kafka-l0 kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic orders --from-beginning