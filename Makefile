BINARY = coco

# Build binary
build: install tidy linux
generate-ts: 
	export PATH=$PATH:$(go env GOPATH)/bin;
	swag init --pd --parseDepth 1 -g ./cli/main.go;
	npx swagger-typescript-api -p ./docs/swagger.json -o ./sign-in-ui/src -n api.ts;
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
node:
	cd sign-in-ui;
	yarn build;