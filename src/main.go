package main

import (
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"taurus/infrastructures"
	proto "taurus/infrastructures/proto"
	"taurus/injectors"
)

// kelseyhightower/envconfig とか使ってもよいがDBのホストとかファイルに書くのも微妙なので・・・
type Env struct {
	DatabaseHost string
	DatabasePort string
	DatabaseUser string
	DatabasePass string
	DatabaseName string
	ListenPort   string
}

func main() {
	log.Println("Initializing taurus...")

	// 環境変数の取得
	env, envErr := getEnvironment()
	if envErr != nil {
		// 何も出来ないのでpanic
		log.Panicf("Failed get environment: %v", envErr)
		return
	}

	// listenポートの取得
	listen, err := net.Listen("tcp", env.ListenPort)
	if err != nil {
		// 何も出来ないのでpanic
		log.Panicf("Port conflict detected with port %s", env.ListenPort)
		return
	}

	// DBのコネクション確立とコネクションプールの生成
	database, dbErr := infrastructures.NewDatabase(env.DatabaseHost, env.DatabasePort, env.DatabaseUser, env.DatabasePass, env.DatabaseName)
	if dbErr != nil {
		// 何も出来ないのでpanic
		log.Panicf("Failed connect to database: %v", dbErr)
		return
	}

	// grpcサーバーインスタンスの生成
	server := infrastructures.CreateServer()

	// サービスのインスタンス生成
	service := injectors.InjectUserService(database)

	// grpcサーバーインスタンスにserviceを登録する
	proto.RegisterUserServiceServer(server, service)

	// お片付け
	defer func() {
		if err := database.Close(); err != nil {
			// 自分で消してください...
			log.Printf("Failed close database connection: %v", err)
		}
	}()

	// ここごと殺されても困るので別ルーチンでサーバーを起動する
	go func() {
		log.Printf("Start taurus port: %s", env.ListenPort)
		if err := server.Serve(listen); err != nil {
			// 自分で殺してください...
			log.Printf("Failed run taurus: %v", err)
		}
	}()

	// OSからキルシグナルが来たときの処理
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Stopping taurus...")
	server.GracefulStop()
}

func getEnvironment() (*Env, error) {
	dbHost, dbHostOk := os.LookupEnv("DATABASE_URL")
	log.Printf("dbHost: %s, dbHostOk: %t", dbHost, dbHostOk)
	if !dbHostOk {
		return nil, errors.New("DATABASE_URL is not found")
	}
	dbPort, dbPortOk := os.LookupEnv("DATABASE_PORT")
	if !dbPortOk {
		return nil, errors.New("DATABASE_PORT is not found")
	}
	dbUser, dbUserOk := os.LookupEnv("DATABASE_USER")
	if !dbUserOk {
		return nil, errors.New("DATABASE_USER is not found")
	}
	dbPass, dbPassOk := os.LookupEnv("DATABASE_PASSWORD")
	if !dbPassOk {
		return nil, errors.New("DATABASE_PASSWORD is not found")
	}
	dbName, dbNameOk := os.LookupEnv("DATABASE_NAME")
	if !dbNameOk {
		return nil, errors.New("DATABASE_NAME is not found")
	}
	listenPort, listenPortOk := os.LookupEnv("LISTEN_PORT")
	if !listenPortOk {
		return nil, errors.New("LISTEN_PORT is not found")
	}
	return &Env{
		DatabaseHost: dbHost,
		DatabasePort: dbPort,
		DatabaseUser: dbUser,
		DatabasePass: dbPass,
		DatabaseName: dbName,
		ListenPort:   listenPort,
	}, nil
}
