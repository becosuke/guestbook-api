PROJECT_REPOSITORY=github.com/becosuke/guestbook-api
PROJECT_NAME=guestbook-api
DOCKER_NAME=guestbook-golang
DOCKER_GOPATH=/go
OUTPUT=output
TIMESTAMP=$(shell date +%Y%m%d%H%M%S)
PRECOMMAND=docker exec -e 'CGO_ENABLED=0' -e 'GOOS=linux' ${DOCKER_NAME}

all: mod test build

clean: clean-build clean-test

clean-build:
	${PRECOMMAND} go clean

clean-test:
	${PRECOMMAND} go clean -testcache

mod:
	${PRECOMMAND} go mod vendor

mod-tidy:
	${PRECOMMAND} go mod tidy

fmt:
	${PRECOMMAND} sh -c "find . -type d \\( -name .git -o -name vendor -o -name output \\) -prune -o -type f -name *.go -print | xargs -n1 go fmt"

vet:
	${PRECOMMAND} sh -c "find . -type d \\( -name .git -o -name vendor -o -name output \\) -prune -o -type d -print | xargs -IXXX sh -c 'find XXX -maxdepth 1 -type f -name *.go -print | xargs --no-run-if-empty go vet' || :"

lint:
	${PRECOMMAND} sh -c "find . -type d \\( -name .git -o -name vendor -o -name output \\) -prune -o -type f -name *.go -print | xargs -n1 golint | grep -v 'should have comment' || :"

test:
	${PRECOMMAND} go test -v ${PROJECT_REPOSITORY}/application/controller

build:
	${PRECOMMAND} go build -a -installsuffix cgo -ldflags '-w' -o ${OUTPUT}/main main.go
	cd ${OUTPUT} && docker build . --no-cache -t ${PROJECT_NAME}:latest -t ${PROJECT_NAME}:${TIMESTAMP}
