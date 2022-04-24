.PHONY: all test update gen fontend proto build_dir frontend_settings


all: build_dir frontend
	go mod tidy && go build -v -o ./build/pns ./cmd/mono_pns_backend/main.go && cp -R ./config ./build

build_dir:
	mkdir -p ./build

frontend:
	GOARCH=wasm GOOS=js go build -v -o ./web/app.wasm ./cmd/pns_frontend/main.go && cp -R ./web ./build/web

test:
	go test -v ./...

gen:
	gf gen dao -g mysql -n -s -l "mysql:root:pns_root@tcp(127.0.0.1:3306)/pns"

frontend_settings:
	bash ./scripts/generate_settings.sh ./web/settings.json ./internal/admin/frontend/settings/raw.go settings

proto:
	cd proto && bash generate.sh