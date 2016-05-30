all:
	cd ./src && go build -o ../uct

install:
	cp ./uct /usr/bin/uct
