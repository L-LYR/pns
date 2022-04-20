.PHONY: android windows mobile_settings linux

mobile_settings:
	bash ../../scripts/generate_settings.sh ./config/settings.json ./raw.go main

android: mobile_settings
	mkdir -p build && fyne package --os android --appID pns.mobile.demo --icon ./img/logo.png --name PNS-Mobile && mv *.apk ./build

windows:
	go build -o ./build/pns.exe -ldflags -H=windowsgui .

linux:
	go build -o ./build/pns .