package main

import (
	"github.com/go-pg/pg"
)

//House ...
type House struct {
	Text       string   `pg:"-"`
	tableName  struct{} `pg:"fias.house"`
	HOUSEID    string   `pg:",pk,type:uuid"`
	HOUSEGUID  string   `pg:"type:uuid"`
	AOGUID     string   `pg:"type:uuid"`
	HOUSENUM   string
	STRSTATUS  string `pg:"type:integer"`
	ESTSTATUS  string `pg:"type:integer"`
	STATSTATUS string `pg:"type:integer"`
	IFNSFL     string
	IFNSUL     string
	TERRIFNSFL string
	TERRIFNSUL string
	OKATO      string
	OKTMO      string
	POSTALCODE string
	STARTDATE  string `pg:"type:date"`
	ENDDATE    string `pg:"type:date"`
	UPDATEDATE string `pg:"type:date"`
	COUNTER    string `pg:"type:integer"`
	DIVTYPE    string `pg:"type:integer"`
	REGIONCODE string
	NORMDOC    string `pg:"type:uuid"`
	BUILDNUM   string
	CADNUM     string
	STRUCNUM   string
}

//CreateTable ...
func (h *House) CreateTable(db *pg.DB) {
	createTable(h, db)
}

//DropTable ...
func (h *House) DropTable(db *pg.DB) {
	dropTable(h, db)
}

//DropNotRelevantRows ...
func (h *House) DropNotRelevantRows(db *pg.DB) {
	//TODO: Определить, есть ли что удалять
	return
}

//CreateIndexes ...
func (h *House) CreateIndexes(db *pg.DB) {
	//TODO: Определить, нужны ли индексы
	return
}

//Parse ...
func (h *House) Parse(folder string, db *pg.DB) {
	parseFolder(folder, h, db)
}
