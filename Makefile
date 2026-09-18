BINARY    := osmonitor
IMAGE     := osmonitor:latest
# Architectures for the multi-arch targets. Override: make docker-buildx PLATFORMS=linux/arm64
PLATFORMS := linux/amd64,linux/arm64

# Pure-Go build. Avoids needing a C toolchain (and the Xcode license on macOS).
export CGO_ENABLED := 0

# Scoped: ./... would descend into ui/node_modules, which contains a Go package.
PKGS := ./cmd/... ./internal/... ./test/...

# Where the compiled dashboard is embedded from (see internal/web/web.go).
DIST := internal/web/dist

.PHONY: run build test tidy ui-install ui-dev ui-build ui-clean docker-build docker-buildx docker-push docker-run docker-up docker-down

run:
	go run ./cmd/osmonitor

build:
	go build -trimpath -ldflags="-s -w" -o bin/$(BINARY) ./cmd/osmonitor

test:
	go vet $(PKGS)
	go test $(PKGS)

tidy:
	go mod tidy

ui-install:
	cd ui && npm install

ui-dev:
	cd ui && npm run dev

# Generate the static dashboard and copy it into the Go embed directory.
# Run this before `make build` to get a binary that serves the UI.
ui-build:
	cd ui && npm run generate
	find $(DIST) -mindepth 1 -not -name .gitkeep -delete
	cp -R ui/.output/public/. $(DIST)/

ui-clean:
	find $(DIST) -mindepth 1 -not -name .gitkeep -delete

docker-build:
	docker build -t $(IMAGE) .

# Multi-arch image loaded into the local daemon (needs the containerd image
# store, which Docker Desktop enables by default).
docker-buildx:
	docker buildx build --platform $(PLATFORMS) -t $(IMAGE) --load .

# Multi-arch image pushed to a registry, e.g.
#   make docker-push IMAGE=ghcr.io/you/osmonitor:1.0.0
docker-push:
	docker buildx build --platform $(PLATFORMS) -t $(IMAGE) --push .

docker-run: docker-build
	docker run --rm -p 8080:8080 --pid=host \
		-v /proc:/host/proc:ro -v /sys:/host/sys:ro -v /etc:/host/etc:ro \
		$(IMAGE)

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down
