compose_api:
	docker compose -f docker/api.yml up --build -d

compose_debug:
	docker compose -f docker/api-debug.yml up --build -d

doc_gen:
	cd backend && swag init -g pkg/handler/v1/handler.go

debug_api:
	make doc_gen && \
	cd backend && go build -o ../companion.bin ./cmd/main.go && cd .. && \
	DB_URL=postgresql://postgres:companion@localhost:5000/companionai \
	LLM_URL=http://localhost:11434 \
	API_PORT=8000 \
	JWT_AUTH_METHOD=HS256 \
	 ./companion.bin