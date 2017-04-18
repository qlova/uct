all:
	cd ./src && go-bindata -pkg "uct" -o ./data.go ./data/
	cd ./src/uct && go build -o ../../uct

install:
	cp ./uct /usr/bin/uct
