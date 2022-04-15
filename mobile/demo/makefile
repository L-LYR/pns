.PHONY: all

mobile_settings:
	bash ../../scripts/generate_settings.sh ./config/settings.json ./raw.go main

all: mobile_settings
	mkdir -p build && fyne package --os android --appID pns.mobile.demo --icon ./img/logo.png --name PNS-Mobile && mv *.apk ./build
