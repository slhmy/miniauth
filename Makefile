.PHONY: swag
swag:
	swag fmt
	swag init --ot json -o api
	openapi-generator-cli generate -i /local/api/swagger.json -g typescript-axios -o /local/website/src/api

.PHONY: build-frontend
build-frontend:
	cd website && pnpm install && pnpm build

.PHONY: build
build: build-frontend
	go build -o bin/ .

.PHONY: build-backend
build-backend:
	go build -o bin/ .

# Docker targets
.PHONY: docker-build
docker-build:
	docker build -t miniauth:latest .

.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 -v miniauth_data:/data miniauth:latest
