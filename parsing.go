package main

import (
	"io/ioutil"
	"reflect"
	"strings"
	"sync"

	"go.uber.org/zap"

	"github.com/LindsayBradford/go-dbf/godbf"
	"github.com/go-pg/pg"
)

func parseFolder(folder string, objType interface{}, db *pg.DB) {
	var wg = &sync.WaitGroup{}
	var ch = make(chan string, 100)
	var filename string

	switch t := objType.(type) {
	case *AddressObject:
		filename = "ADDROB"
	case *House:
		filename = "HOUSE"
	case *Socrbase:
		filename = "SOCRBASE"
	default:
		logger.Error("Type not found:",
			zap.String("type:", t.(string)),
		)
	}

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		logger.Error("Folder read error:",
			zap.String("folder:", folder),
			zap.String("error:", err.Error()),
		)
	}

	wg.Add(config.DBF.Threads)
	for i := 0; i < config.DBF.Threads; i++ {
		go parseDBF(objType, ch, wg, db)
	}

	for _, f := range files {
		if strings.Index(f.Name(), filename) == 0 {
			ch <- folder + f.Name()
		}
	}
	close(ch)
	wg.Wait()
}

func parseDBF(objType interface{}, in chan string, wg *sync.WaitGroup, db *pg.DB) {
	var wgInsert = &sync.WaitGroup{}
	var chInsert = make(chan interface{}, 1000)

	insertWorkersMaxCnt := getMaxInsertWorkers(db)

	wgInsert.Add(insertWorkersMaxCnt)
	for i := 0; i < insertWorkersMaxCnt; i++ {
		if config.Program.Update {
			go updateAsync(chInsert, wgInsert, db)
		} else {
			go insertAsync(chInsert, wgInsert, db)
		}
	}

	go func() {
		for {
			path, ok := <-in
			if !ok {
				wg.Done()
				return
			}
			if len(path) != 0 {
				logger.Info("Parse dbf",
					zap.String("stage:", "start"),
					zap.String("path:", path),
				)

				dbfTable, err := godbf.NewFromFile(path, config.DBF.Codepage)
				if err != nil {
					logger.Error("DBF-file read error:",
						zap.String("path:", path),
						zap.String("codepage:", config.DBF.Codepage),
						zap.String("error:", err.Error()),
					)
				}

				var a DBObject

				for i := 0; i < dbfTable.NumberOfRecords(); i++ {

					switch t := objType.(type) {
					case *AddressObject:
						a = &AddressObject{}
					case *House:
						a = &House{}
					case *Socrbase:
						a = &Socrbase{}
					default:
						logger.Error("Type not found:",
							zap.String("type:", t.(string)),
						)
					}

					s := reflect.ValueOf(a).Elem()
					for _, field := range dbfTable.Fields() {
						val, _ := dbfTable.FieldValueByName(i, field.Name())
						if s.Kind() == reflect.Struct {
							f := s.FieldByName(field.Name())
							if f.IsValid() {
								if f.CanSet() {
									f.SetString(val)
								}
							}
						}
					}
					chInsert <- a
				}
				logger.Info("Parse dbf",
					zap.String("stage:", "finish"),
					zap.String("path:", path),
					zap.Int("rows in file:", dbfTable.NumberOfRecords()),
				)
			}
		}
	}()
}
