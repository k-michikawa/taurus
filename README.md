# taurus

### Description

Go で gorm 使って gRPC で CRUD するやつ

Go 1.15.6

検証用に grpcurl 入れておくと良いかも

mac: `brew install grpcurl`

example:

```sh
$ grpcurl -plaintext \
          -import-path ./proto \
          -proto user.proto \
          -d '{"name": "taurus", "email": "taurus@test.com", "password": "password"}' \
          localhost:9010 \
          leo.UserService/PostUser
```

### tools

golang-migrate のインストール

```sh
$ brew install golang-migrate
```

migration ファイルの生成

```sh
$ migrate create -ext sql -dir migrations -seq /* filename */
```

protobuf のインストール

```sh
$ brew install protobuf
```

protodep のインストール

```sh
$ go get github.com/stormcat24/protodep
```

OR

```sh
$ wget https://github.com/stormcat24/protodep/releases/download/0.0.8/protodep_darwin_amd64.tar.gz
$ tar -xf protodep_darwin_amd64.tar.gz
$ mv protodep /usr/local/bin/
```

AFTER

```sh
$ ssh-add ~/.ssh/id_rsa
```

### run

起動するまで

1. goenv とかで Go の環境作っておく
2. protodep 落としてくる
3. `$ protodep up`
4. `$ protoc --proto_path ./proto --go_out=plugins=grpc:src/infrastructures/proto user.proto`
5. `$ docker-compose up -d taurus-db`
6. `$ (cd ./migrations && go install && go run .)`
7. `$ docker stop taurus-db`
8. `$ docker-compose up`
