.PHONY: command clean

GRAPHQL_CMD=protoc-gen-graphql
VERSION=$(or ${tag}, dev)

command: plugin clean
	cd ${GRAPHQL_CMD} && \
		go build \
			-ldflags "-X main.version=${VERSION}" \
			-o ../dist/${GRAPHQL_CMD}

lint:
	golangci-lint run

plugin:
	protoc -I $(shell brew --prefix protobuf)/include/google \
		-I include/graphql \
		--go_out=./graphql \
		include/graphql/graphql.proto
	mv graphql/github.com/rleszilm/grpc-graphql-gateway/graphql/graphql.pb.go graphql/
	rm -rf graphql/github.com

test:
	go list ./... | xargs go test

build: test
	protoc -I google \
		-I include/graphql \
		--go_out=./graphql \
		include/graphql/graphql.proto
	mv graphql/github.com/rleszilm/grpc-graphql-gateway/graphql/graphql.pb.go graphql/
	rm -rf graphql/github.com

clean:
	rm -rf ./dist/*

all: clean build
	cd ${GRAPHQL_CMD} && GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o ../dist/${GRAPHQL_CMD}.darwin
	cd ${GRAPHQL_CMD} && GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o ../dist/${GRAPHQL_CMD}.linux
