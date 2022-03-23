all: frontend
	mkdir -p ./build && go mod tidy && go build -v -o ./build/pns ./cmd/mono_pns_backend/main.go

frontend:
	GOARCH=wasm GOOS=js go build -v -o ./web/app.wasm ./cmd/pns_frontend/main.go

test:
	go test -v ./...

update:
	bash ./scripts/deploy.sh update

gen:
	gf gen dao -g pns -n -s -l "mysql:root:pns_root@tcp(127.0.0.1:3306)/pns"

.PHONY:
	all test update gen
