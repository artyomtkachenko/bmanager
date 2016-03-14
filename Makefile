export PATH := /usr/local/go/bin:$(PATH)
export GOPATH := $(PWD)
export TMPDIR := $(PWD)/tmp
DEPS = $(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
VETARGS=-asmdecl -atomic -bool -buildtags -copylocks -methods \
									-nilfunc -printf -rangeloops -shift -structtags -unsafeptr

VERSION=$(shell [ -n "${GO_PIPELINE_COUNTER}" ] && echo "${GO_PIPELINE_COUNTER}" || echo "dev" )

all: test

deps:
		mkdir -p tmp
		go get -v -u golang.org/x/net/html

fmt:
		mkdir -p tmp
		go fmt ./...

vet: fmt
		go tool vet ${VETARGS} src/balancer-manager/
		go tool vet ${VETARGS} main.go

test: vet

		go test -v balancer-manager/apache
		go test -v .

cover: test
		go test -coverprofile=cover.out  balancer-manager/apache/
		go tool cover -html=cover.out

bench: test
		go test balancer-manager/apache -bench=. -cpuprofile=cpu.pprof
		go tool pprof --pdf apache.test cpu.pprof >cpu_stats.pdf
		open cpu_stats.pdf

build: test
		go build -ldflags "-X main.version=${VERSION}" -o bin/balancer-manager main.go
		rm -rf tmp

.PHONY: all deps vet test cover bench build fmt

