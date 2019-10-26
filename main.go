package main

import (
	"time"

	"go.uber.org/zap"
)

const configPath = `config.yml`

var logger *zap.Logger
var config Config

func main() {
	started := time.Now()
	getConfig(configPath)
	logger, _ = NewLogger()

	db := conn()
	defer db.Close()

	logger.Info("Success connection.",
		zap.Int("max_connections:", getMaxConnections(db)),
		zap.Int("max insert worker:", getMaxInsertWorkers(db)),
	)

	addExtensionPgTrgm(db)

	dbobjects := []DBObject{
		&AddressObject{},
		&House{},
		&Socrbase{},
	}

	for _, obj := range dbobjects {
		if config.Program.DropTables && !config.Program.Update {
			obj.DropTable(db)
		}

		//TODO: Проверить на наличие таблиц. Выдать ошибку и выйти если хотя бы одной таблицы нет
		if config.Program.CreateTables && !config.Program.Update {
			obj.CreateTable(db)
			obj.CreateIndexes(db)
		}

		obj.Parse(config.DBF.Folder, db)

		if config.Program.DropIrrelevantRows {
			obj.DropNotRelevantRows(db)
		}
	}

	time.Sleep(5 * time.Second)

	logger.Info("Work ended.",
		zap.Duration("worked time:", time.Since(started)))
}

//TODO: Настройка логгера через .yml
func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		config.Log.Path,
		"stdout",
	}
	return cfg.Build()
}

//TODO: Прогресс-бар в консоле
//TODO: Автоматическое обновление с сайта fias.nalog.ru
