package main

/**
 * CLIでよしなにやるのめんどくさかったのでいい感じにマイグレートする子
 */

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	// pwdのpath取得
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed get current directory path: %v", err)
	}
	// databaseインスタンスの生成(.Openとか名前ついてるが実際pingすら飛ばしてないっぽい)
	db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/taurus-db?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed connect DB: %v", err)
	}
	// postgresインスタンスの生成(ここでコネクションの確立を行う)
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed generate postgres instance: %v", err)
	}
	// migrateファイルの取得
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s/sql", pwd), "postgres", driver)
	if err != nil {
		log.Fatalf("Failed generate migrate instance: %v", err)
	}
	if err := m.Up(); err != nil {
		log.Printf("Failed migration: %v", err)
	}
}
