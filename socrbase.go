package main

import (
	"github.com/go-pg/pg"
)

type Socrbase struct {
	tableName struct{} `pg:"fias.socrbase"`
	LEVEL     string
	SCNAME    string
	SOCRNAME  string
	KOD_T_ST  string `pg:",pk"`
}

func (s *Socrbase) CreateTable(db *pg.DB) {
	createTable(s, db)
}

func (s *Socrbase) DropTable(db *pg.DB) {
	dropTable(s, db)
}

func (s *Socrbase) DropNotRelevantRows(db *pg.DB) {
	//TODO: Определить, есть ли что удалять
	return
}

func (s *Socrbase) CreateIndexes(db *pg.DB) {
	//TODO: Определить, нужны ли индексы
	return
}

func (s *Socrbase) Parse(folder string, db *pg.DB) {
	parseFolder(folder, s, db)
}
