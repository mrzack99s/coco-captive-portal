BINARY = coco

# Build binary
build: install tidy linux
generate-ts: 
	export PATH=$PATH:$(go env GOPATH)/bin;
	swag init --pd --parseDepth 1 -g ./cli/main.go;
	npx swagger-typescript-api -p ./docs/swagger.json -o ./sign-in-dev/src -n api.ts;
	cp ./sign-in-dev/src/api.ts ./operator-dev/src/api.ts
install-dev-tools:
	sudo snap install go --channel=1.17/stable --classic;
	sudo snap install node --classic;
	sudo snap install redis;
	sudo apt-get update;
	sudo apt install -y build-essential libpcap0.8-dev;
	curl -o- -L https://yarnpkg.com/install.sh | sudo bash
	go install github.com/swaggo/swag/cmd/swag@latest;
install:
	go install ./cli/main.go;
tidy:
	go mod tidy;
linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ${BINARY} ./cli/main.go ;
	-chmod +x ${BINARY};
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY}-installer ./installer/cli.go ;
	-chmod +x ${BINARY}-installer;
build-node:
	cd sign-in-dev;
	yarn build;
	mv build ../dist-auth-ui
	cd ../operator-dev;
	yarn build;
	mv build ../dist-operator-ui