COMPONENT_NAME:=canned
VERSION=$(shell cat $(CURDIR)/version)
DOCKER_TAG=$(COMPONENT_NAME):$(VERSION)

GO111MODULE=on
export GO111MODULE
GOFLAGS+=-mod=vendor
export GOFLAGS

LDFLAGS=-ldflags "-w -s -X main.version=${COMPONENT_VERSION}"

.PHONY: all build docker format lint vendor clean

all: build

build:
	$(info Building $(COMPONENT_NAME))
	@go build ${LDFLAGS} -o ${COMPONENT_NAME}

docker: build
	$(info Builder docker image for $(COMPONENT_NAME))
	@docker build . -t $(DOCKER_TAG)

docker_install:
	$(info Installing $(COMPONENT_NAME))
	@CGO_ENABLED=0 go install ${LDFLAGS}

test:
	$(info Performing tests)
	@go test -v -cover .

coverage:
	$(info Performing coverage)
	@go test -coverprofile=coverage.out && go tool cover -func=coverage.out && go tool cover -html=coverage.out

format:
	$(info Formatting code)
	go fmt ./...

lint:
	$(info Performing static analysis)
	@golint .

vendor:
	@$(info Updating vendored modules)
	@go mod tidy && go mod vendor

clean:
	$(info Cleaning $(COMPONENT_NAME))
	@$(RM) ${COMPONENT_NAME}
	$(info Removing test artifacts...)
	@$(RM) coverage.*


