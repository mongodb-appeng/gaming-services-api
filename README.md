# Gaming Platform Services API

__NOTE__ This project structure follows the layout specified [here](https://github.com/golang-standards/project-layout)

## TODO

- OpenAPI spec
- add hosting
- add testing
- dynamic CORS support
- See other notes in trello board

## Set up

Note: Go must be installed and configured. See [docs](https://golang.org/doc/install)

```bash
mkdir -p $GOPATH/src/github.com/desteves/
cd $GOPATH/src/github.com/desteves/
git clone https://github.com/desteves/babysteps.git
cd babysteps
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
   -v  "$PWD":/go/src/github.com/desteves/babysteps \
   -p 8888:8888  \
   --hostname gps  \
   game-docker

## once inside the container
cd /go/src/github.com/desteves/babysteps &&  go run cmd/gameplatformservices/main.go -c configs/config.yaml
```
