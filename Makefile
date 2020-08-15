APPLICATION := $(shell basename `pwd`)
BUILD_RFC3339 := $(shell date -u +"%Y-%m-%dT%H:%M:%S+00:00")
PACKAGE := $(shell git remote get-url --push origin | sed -E 's/.+[@|/](.+)\.(.+).git/\1.\2/' | sed 's/\:/\//')
REVISION := $(shell git rev-parse HEAD)
VERSION := $(shell git describe --tags)
DESCRIPTION := $(shell curl -s https://api.github.com/repos/${PACKAGE} \
    | grep '"description".*' \
    | head -n 1 \
    | cut -d '"' -f 4)
WORKDIR := $(shell pwd)

GO_LDFLAGS := "-w -s \
	-X github.com/jnovack/release.Application=${APPLICATION} \
	-X github.com/jnovack/release.BuildRFC3339=${BUILD_RFC3339} \
	-X github.com/jnovack/release.Package=${PACKAGE} \
	-X github.com/jnovack/release.Revision=${REVISION} \
	-X github.com/jnovack/release.Version=${VERSION} \
	"

DOCKER_BUILD_ARGS := \
	--build-arg APPLICATION=${APPLICATION} \
	--build-arg BUILD_RFC3339=v${BUILD_RFC3339} \
	--build-arg DESCRIPTION="${DESCRIPTION}" \
	--build-arg PACKAGE=${PACKAGE} \
	--build-arg REVISION=${REVISION} \
	--build-arg VERSION=${VERSION} \
	--progress auto

all: build

.PHONY: assets
assets:
	cd assets/ && curl -LO https://raw.githubusercontent.com/wp-statistics/GeoLite2-City/master/GeoLite2-City.mmdb.gz
	cd assets/ && curl -LO https://raw.githubusercontent.com/wp-statistics/GeoLite2-Country/master/GeoLite2-Country.mmdb.gz
	# cd assets/ && curl -LO https://github.com/wp-statistics/GeoLite2-ASN/raw/master/GeoLite2-ASN.mmdb.gz
	gunzip -fq assets/*.gz

.PHONY: build
build:
	go build -o bin/${APPLICATION} -ldflags $(GO_LDFLAGS) cmd/*/*

.PHONY: docker
docker:
	docker build ${DOCKER_BUILD_ARGS} -t ${APPLICATION}:latest .

.PHONY: test
test:
	go test ./... -timeout 30s