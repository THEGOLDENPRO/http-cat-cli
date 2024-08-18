build:
	go build

install:
	sudo cp ./http-cat-cli /usr/bin/http-cat

uninstall:
	sudo rm /usr/bin/http-cat