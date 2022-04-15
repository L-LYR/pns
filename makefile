.PHONY: all test update gen mobile fontend proto

all: frontend
	mkdir -p ./build && go mod tidy && go build -v -o ./build/pns ./cmd/mono_pns_backend/main.go && mv -r ./config ./build

proto:
	cd proto && bash generate.sh

frontend_settings:
	bash ./scripts/generate_settings.sh ./web/settings.json ./internal/admin/frontend/settings/raw.go settings

frontend: frontend_settings
	GOARCH=wasm GOOS=js go build -v -o ./web/app.wasm ./cmd/pns_frontend/main.go && mv -r ./web ./build

mobile:
	cd ./mobile/demo && make all

test:
	go test -v ./...

update:
	bash ./scripts/deploy.sh update

gen:
	gf gen dao -g mysql -n -s -l "mysql:root:pns_root@tcp(127.0.0.1:3306)/pns"
