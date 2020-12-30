# taurus

Go で google.golang.org/grpc と gorm を使って gRPC で CRUD するやつ

検証用に grpcurl 入れておくと良いかも

mac: `brew install grpcurl`

example: `grpcurl -plaintext -import-path ./proto -proto user.proto -d '{"name": "taurus", "email": "taurus@test.com", "password": "password"}' localhost:9010 taurus.UserService/PostUser`

migration ファイルの生成

1. install golang-migrate `brew install golang-migrate`
2. `migrate create -ext sql -dir migrations -seq /* table-name */`

migrate 方法

1. `docker-compose up`で DB のコンテナを立ち上げる(多分初回は migrate してないので grpc の方は落ちる)
2. `go run ./migrations`
