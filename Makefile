# Perintah harian developer Fablerace.

.PHONY: build run vet test tidy up down logs next relasi

build:
	go build ./...

run:
	go run ./cmd/api

vet:
	go vet ./...

test:
	go test ./...

tidy:
	go mod tidy

up:
	docker compose -f deploy/docker-compose.yml up --build -d

down:
	docker compose -f deploy/docker-compose.yml down

logs:
	docker compose -f deploy/docker-compose.yml logs -f api

# Backlog teknis yang hidup di kode (docs/lampiran-b-konvensi-kode.md §4).
next:
	@grep -rn "NEXT:" internal/ cmd/ || true

relasi:
	@grep -rn "RELASI:" internal/ cmd/ || true
