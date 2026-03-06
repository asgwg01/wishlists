include .env
export

postgres-migrate-up:
# authService
	migrate -path authService/migrations -database "postgres://${AUTH_DB_USER}:${AUTH_DB_PASSWORD}@localhost:${AUTH_DB_PORT}/${AUTH_DB_NAME}?sslmode=disable" up
# authService
	migrate -path wishlistService/migrations -database "postgres://${WISHLIST_DB_USER}:${WISHLIST_DB_PASSWORD}@localhost:${WISHLIST_DB_PORT}/${WISHLIST_DB_NAME}?sslmode=disable" up

postgres-migrate-down:
# authService
	migrate -path authService/migrations -database "postgres://${AUTH_DB_USER}:${AUTH_DB_PASSWORD}@localhost:${AUTH_DB_PORT}/${AUTH_DB_NAME}?sslmode=disable" down
# authService
	migrate -path wishlistService/migrations -database "postgres://${WISHLIST_DB_USER}:${WISHLIST_DB_PASSWORD}@localhost:${WISHLIST_DB_PORT}/${WISHLIST_DB_NAME}?sslmode=disable" down

# Запуск инфраструктуры
wishlist-infra-deploy:
	mkdir -p .volumes/zookeeper/data
	mkdir -p .volumes/zookeeper/log
	mkdir -p .volumes/kafka/data
	mkdir -p .volumes/pg_service_auth
	mkdir -p .volumes/pg_service_wishlist
	mkdir -p .volumes/pg_service_notify
	mkdir -p .volumes/logbull
	mkdir -p .volumes/redis
	docker compose up --build -d pg_service_auth pg_service_wishlist pg_service_notify redis_service zookeeper kafka kafka_ui logbull

wishlist-infra-undeploy:
	docker compose down pg_service_auth pg_service_wishlist pg_service_notify redis_service zookeeper kafka kafka_ui logbull

wishlist-deploy:
	mkdir -p .volumes/zookeeper/data
	mkdir -p .volumes/zookeeper/log
	mkdir -p .volumes/kafka/data
	mkdir -p .volumes/pg_service_auth
	mkdir -p .volumes/pg_service_wishlist
	mkdir -p .volumes/pg_service_notify
	mkdir -p .volumes/logbull
	mkdir -p .volumes/redis
	docker compose up --build -d

# Остановка инфраструктуры
wishlist-undeploy:
	docker compose down

# Генерация protobuf файлов
wishlist-protobuf-generate:
	@echo "Generating protobuf files..."
	protoc --go_out=wishlistService/api --go_opt=paths=source_relative \
		--go-grpc_out=wishlistService/api --go-grpc_opt=paths=source_relative \
		wishlistService/api/proto/*.proto



# Запуск конкретного сервиса локально
# authService
wishlist-authService-run:
	cd authService && go run cmd/authService/main.go

# wishlistService
wishlist-wishlistService-run:
	cd wishlistService && go run cmd/wishlistService/main.go

# apiGatewayService
wishlist-apiGateway-run:
	cd apiGateway && go run cmd/apiGateway/main.go

# notificationService
wishlist-notificationService-run:
	cd notificationService && go run cmd/notificationService/main.go

# webClient
wishlist-webClient-run:
	cd webClient && go run cmd/webClient/main.go

wishlist-proto-gen:
	@echo "Generating protobuf code..."
	protoc \
		--go_out=. \
		--go_opt=module=wishlistApp \
		--go-grpc_out=. \
		--go-grpc_opt=module=wishlistApp \
		pkg/proto/**/v1/*.proto
	@echo "Done!"

# Очистка сгенерированного кода
wishlist-proto-clean:
	@echo "Cleaning generated protobuf code..."
	cd ./pkg && find . -name "*.pb.go" -delete
	cd ./pkg && find . -name "*_grpc.pb.go" -delete
	@echo "Done!"

# Перегенерация всего
wishlist-proto-regen: wishlist-proto-clean wishlist-proto-gen


swagger-init:
	cd apiGateway && swag init -g cmd/apiGateway/main.go

# test:
# 	go test TODO!!!