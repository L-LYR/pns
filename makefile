all:
	mkdir -p ./build && go mod tidy && go build -v -o ./build/pns ./app/mono_pns/main.go

test:
	go test -v ./...

update:
	bash ./scripts/deploy.sh update

.PHONY:
	all test update
