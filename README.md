# Gaming Platform Services API

__NOTE__ This project structure follows the layout specified [here](https://github.com/golang-standards/project-layout)

## TODO

- OpenAPI spec
- add static hosting
- add testing
- dynamic CORS support
- See other notes in trello board

## Set up

Note: Go must be installed and configured. See [docs](https://golang.org/doc/install)

```bash
mkdir -p $GOPATH/src/github.com/mongodb-appeng/
cd $GOPATH/src/github.com/mongodb-appeng/
git clone https://github.com/mongodb-appeng/gaming-services-api.git
cd gaming-services-api
```

## Localhost Run

Note: Navigate to the root dir of this project

```bash
go mod
go get -d -v ./...
go install -v ./...
go run cmd/gameplatformservices/main.go -c configs/config.yaml
```

## Dockerized Run

Note: Navigate to the root dir of this project

```bash
docker build -t game-docker build/package/
docker run --rm  -it \
   -v  "$PWD":/go/src/github.com/mongodb-appeng/gaming-services-api \
   -p 8888:8888  \
   --hostname gps  \
   game-docker

## once inside the container
cd /go/src/github.com/mongodb-appeng/gaming-services-api &&  go run cmd/gameplatformservices/main.go -c configs/config.yaml
```
