.PHONY: all test update gen mobile fontend proto

all: frontend
	mkdir -p ./build && go mod tidy && go build -v -o ./build/pns ./cmd/mono_pns_backend/main.go

proto:
	cd proto && bash generate.sh

frontend:
	GOARCH=wasm GOOS=js go build -v -o ./web/app.wasm ./cmd/pns_frontend/main.go

mobile:
	cd ./cmd/pns_mobile && fyne package --os android --appID my.demo.app --icon ../../mobile/img/logo.png --name pns-mobile && mv *.apk ../../build

test:
	go test -v ./...

update:
	bash ./scripts/deploy.sh update

gen:
	gf gen dao -g mysql -n -s -l "mysql:root:pns_root@tcp(127.0.0.1:3306)/pns"
