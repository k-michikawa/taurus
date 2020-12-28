# taurus

Go で google.golang.org/grpc と gorm を使って gRPC で CRUD するやつ

検証用に grpcurl 入れておくと良いかも

mac: `brew install grpcurl`

example: `grpcurl -plaintext -import-path ./proto -proto user.proto -d '{"name": "taurus", "email": "taurus@test.com", "password": "password"}' localhost:9010 taurus.UserService/PostUser`
