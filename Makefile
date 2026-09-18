BINARY := osmonitor
IMAGE  := osmonitor:latest

# Pure-Go build. Avoids needing a C toolchain (and the Xcode license on macOS).
export CGO_ENABLED := 0

.PHONY: run build test tidy ui-install ui-dev ui-build docker-build docker-run docker-up docker-down

run:
	go run ./cmd/osmonitor

build:
	go build -trimpath -ldflags="-s -w" -o bin/$(BINARY) ./cmd/osmonitor

# Scoped: ./... would descend into ui/node_modules, which contains a Go package.
PKGS := ./cmd/... ./internal/...

test:
	go vet $(PKGS)
	go test $(PKGS)

tidy:
	go mod tidy

ui-install:
	cd ui && npm install

ui-dev:
	cd ui && npm run dev

ui-build:
	cd ui && npm run build

docker-build:
	docker build -t $(IMAGE) .

docker-run: docker-build
	docker run --rm -p 8080:8080 --pid=host \
		-v /proc:/host/proc:ro -v /sys:/host/sys:ro -v /etc:/host/etc:ro \
		$(IMAGE)

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down
