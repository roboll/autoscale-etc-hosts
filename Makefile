###############################################################################
# release tasks for binary and docker releases
###############################################################################
OWNER     := roboll
REPO      := autoscale-etc-hosts

PROJECT   := github.com/$(OWNER)/$(REPO)
VERSION   := $(shell git describe --tags)
IMAGE_TAG := $(OWNER)/$(REPO):$(VERSION)

GOOS     := linux
GOARCH   := amd64
BINARY   := $(REPO)-$(GOOS)-$(GOARCH)

all: build
build: target/$(BINARY) docker-build
release: clean $(PRE_RELEASE) gh-release-$(BINARY) docker-release

.PHONY: clean
clean: ; rm -rf target

###############################################################################
# pre-release - test and validation steps
###############################################################################
PRE_RELEASE := test

.PHONY: test
test: ; go test ./...

###############################################################################
# release artifacts
###############################################################################
target: ; mkdir -p target

target/$(BINARY): target
	docker run --rm \
		-v $(PWD):/go/src/$(PROJECT) -v $(PWD)/target:/target \
		golang /bin/bash -c \
			"CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
			go get $(PROJECT)/... && \
			go build -a -tags netgo \
			-ldflags '-s -w -X main.release=$(VERSION)' \
			-o /target/$(BINARY) $(PROJECT)"

target/%.tar.gz: target %
	@echo Packaging $* to target/$*.tar.gz.
	@tar czf target/$*.tar.gz -C $* .

###############################################################################
# docker releases
###############################################################################
.PHONY: docker-build docker-release

docker-build: target/$(BINARY)
	docker build -t $(IMAGE_TAG) .

docker-release: target/$(BINARY) docker-build
	docker push $(IMAGE_TAG)

###############################################################################
# github-release - upload a binary release to github releases
#
# requirements:
# - the checked out revision be a pushed tag
# - a github api token (GITHUB_TOKEN)
###############################################################################
API    = https://api.github.com/repos/$(OWNER)/$(REPO)
UPLOAD = https://uploads.github.com/repos/$(OWNER)/$(REPO)

.PHONY: create-gh-release gh-release gh-token
create-gh-release: tag clean-repo gh-token
	$(info Creating Github Release)
	@curl -s -XPOST -H "Authorization: token $(GITHUB_TOKEN)" \
		$(API)/releases -d '{ "tag_name": "$(VERSION)" }' > /dev/null

gh-release-%: tag clean-repo gh-token target/% create-gh-release
	$(info Uploading Release Artifact $* to Github)
	@curl -s \
		-H "Authorization: token $(GITHUB_TOKEN)" \
		$(API)/releases/tags/$(VERSION) |\
	python -c "import json,sys;obj=json.load(sys.stdin);print obj['id']" |\
	curl -s -XPOST \
		-H "Authorization: token $(GITHUB_TOKEN)" \
		-H "Content-Type: application/octet-stream" \
		$(UPLOAD)/releases/$$(xargs )/assets?name=$* \
		--data-binary @target/$* > /dev/null

gh-token:
ifndef GITHUB_TOKEN
	$(error $GITHUB_TOKEN not set)
endif

###############################################################################
# utility
###############################################################################
.PHONY: tag clean
tag:  ; @git describe --tags --exact-match HEAD > /dev/null
clean-repo:
	@git diff --exit-code > /dev/null
	@git diff --cached --exit-code > /dev/null
