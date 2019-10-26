package main

import (
	"fmt"
	"sync"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"go.uber.org/zap"
)

//DBObject ...
type DBObject interface {
	CreateTable(*pg.DB)
	DropTable(*pg.DB)
	DropNotRelevantRows(*pg.DB)
	CreateIndexes(*pg.DB)
	Parse(string, *pg.DB)
}

func conn() *pg.DB {
	return pg.Connect(&pg.Options{
		User:     config.Database.Username,
		Password: config.Database.Password,
		Addr:     config.Database.Address,
		Database: config.Database.DBName,
		//PoolSize: 1000,
	})
}

func createTable(model interface{}, db *pg.DB) {
	err := db.CreateTable(model, &orm.CreateTableOptions{
		Temp: false,
	})

	if err != nil {
		logger.Error("Create table",
			zap.String("table", fmt.Sprintf("%T", model)),
			zap.Error(err),
		)
	}

	logger.Info("Table created",
		zap.String("table", fmt.Sprintf("%T", model)),
	)
}

func dropTable(model interface{}, db *pg.DB) {
	err := db.DropTable(model, &orm.DropTableOptions{})

	if err != nil {
		logger.Error("Drop table",
			zap.String("table", fmt.Sprintf("%T", model)),
			zap.Error(err),
		)
	}

	logger.Info("Table dropped",
		zap.String("table", fmt.Sprintf("%T", model)),
	)
}

func addExtensionPgTrgm(db *pg.DB) {
	_, err := db.Model((*AddressObject)(nil)).Exec(`
			CREATE EXTENSION pg_trgm;			
		`)

	if err != nil {
		logger.Info("Add extension pg_trm",
			zap.String("msg", "extension already exist"),
		)
		return
	}
	logger.Info("Extension successfully added")
}

//TODO: Изменить Insert на Insert or Update для автоматического определения операции
func insertAsync(in chan interface{}, wg *sync.WaitGroup, db *pg.DB) {
	go func() {
		for {
			obj, ok := <-in
			if !ok {
				wg.Done()
				return
			}

			//err := db.Insert(obj)

			_, err := db.Model(obj).
				OnConflict("DO NOTHING").
				Insert()
				// INSERT INTO "books" ("id", "title") VALUES (100, 'my title')
				// ON CONFLICT (id) DO UPDATE SET title = 'title version #1'

			if err != nil {
				logger.Panic("Error insert to database",
					zap.String("record type", fmt.Sprintf("%T", obj)),
					zap.Error(err),
				)
			}
		}
	}()
}

func updateAsync(in chan interface{}, wg *sync.WaitGroup, db *pg.DB) {
	go func() {
		for {
			obj, ok := <-in
			if !ok {
				wg.Done()
				return
			}

			err := db.Update(obj)

			if err != nil {
				logger.Panic("Error update record in database",
					zap.String("record type", fmt.Sprintf("%T", obj)),
					zap.Error(err),
				)
			}
		}
	}()
}

func getMaxConnections(db *pg.DB) int {
	var maxConnections int
	_, err := db.Model((*AddressObject)(nil)).QueryOne(pg.Scan(&maxConnections), `SELECT current_setting('max_connections');`)
	if err != nil {
		logger.Error("Get max_connections",
			zap.Error(err),
		)
	}
	return maxConnections
}

func getMaxInsertWorkers(db *pg.DB) int {
	if getMaxConnections(db) <= 1 {
		return 1
	}
	return getMaxConnections(db)/config.DBF.Threads - 1
}
