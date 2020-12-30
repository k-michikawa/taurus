package infrastructures

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// Databaseに関するアクセスをまとめたやつ以下にもあるが、外からgormに触れてほしくないのでこいつの中にgormを閉じ込めておく
type Database struct {
	// 外からgormに触れて欲しくないのでprivateにしておく
	db *gorm.DB
}

// NewDatabase is sql db constructor
func NewDatabase(host, port, user, password, dbname string) (*Database, error) {
	// 引数からdsnを作る
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Defaultでfalseだが、明示的にやっておく
		// (SELECT... FOR UPDATEとかやらない限りRead以外をトランザクションでラップしてくれる)
		SkipDefaultTransaction: false,
	})
	if err != nil {
		return nil, err
	}

	db.Use(
		// コネクションプールの設定しておく
		dbresolver.Register(dbresolver.Config{}).SetMaxIdleConns(5).SetMaxOpenConns(5),
		// DB関連のミドルウェアほしければここに放り込んでおけば良さそう
	)
	return &Database{db}, nil
}

/* メソッドチェーンするやつ */

// 操作するテーブルや対象のレコード絞るやつ
func (db *Database) SetModel(model interface{}) *Database {
	return &Database{db: db.db.Model(model)}
}

// WHERE するやつ
func (db *Database) Where(query interface{}, args ...interface{}) *Database {
	return &Database{db: db.db.Where(query, args...)}
}

// 否定 WHERE するやつ。gormでは 0, '', false などはクエリに反映されないので <> true とかのクエリを発行する必要がある
func (db *Database) Not(query interface{}, args ...interface{}) *Database {
	return &Database{db: db.db.Not(query, args...)}
}

/* メソッドチェーンしないやつ(実行系) */

// INSERT するやつ。似たものにSaveがあるがSaveはUpdatedAtも自動的に入ってしまう
func (db *Database) Create(value interface{}) error {
	return db.db.Create(value).Error
}

// SELECT WHERE するやつ
func (db *Database) Find(dest interface{}, conds ...interface{}) error {
	return db.db.Find(dest, conds...).Error
}

// LIMIT 1 するやつ
func (db *Database) First(dest interface{}, conds ...interface{}) error {
	return db.db.First(dest, conds...).Error
}

// UPDATE SET WHERE するやつ。引数のmodelとかに主キーとなるもの含んでおけば勝手にWHERE 主キーしてくれる
func (db *Database) Update(model interface{}) error {
	return db.db.Updates(model).Error
}

// DELETE WHERE するやつ。引数のmodelとかに主キーとなるもの含んでおけば勝手にWHERE 主キーしてくれる
// func (db *Database) Delete(model interface{}) error {
// 	return db.db.Delete(model).Error
// }

// gorm.DB に Close メソッドが生えてないのでよしなにやってくれてそうだが一応定義しておく
func (db *Database) Close() error {
	sqlDB, err := db.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
