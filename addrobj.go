package main

import (
	"go.uber.org/zap"

	"github.com/go-pg/pg"
)

//AddressObject Объект адреса
type AddressObject struct {
	Text       string   `pg:"-"`
	tableName  struct{} `pg:"fias.addrobj"`
	AOID       string   `pg:"pk,type:uuid"`
	AOGUID     string   `pg:"type:uuid"`
	FORMALNAME string
	OFFNAME    string
	SHORTNAME  string
	AOLEVEL    string `pg:"type:integer"`
	REGIONCODE string
	AREACODE   string
	AUTOCODE   string
	CITYCODE   string
	CTARCODE   string
	PLACECODE  string
	PLANCODE   string
	STREETCODE string
	EXTRCODE   string
	SEXTCODE   string
	PLAINCODE  string
	CODE       string
	CURRSTATUS string `pg:"type:integer"`
	ACTSTATUS  string `pg:"type:integer"`
	LIVESTATUS string `pg:"type:integer"`
	CENTSTATUS string `pg:"type:integer"`
	OPERSTATUS string `pg:"type:integer"`
	IFNSFL     string
	IFNSUL     string
	OKATO      string
	POSTALCODE string
	STARTDATE  string `pg:"type:date"`
	ENDDATE    string `pg:"type:date"`
	UPDATEDATE string `pg:"type:date"`
	DIVTYPE    string `pg:"type:integer"`
	PARENTGUID string `pg:"type:uuid"`
	NEXTID     string `pg:"type:uuid"`
	OKTMO      string
	PREVID     string `pg:"type:uuid"`
	NORMDOC    string `pg:"type:uuid"`
	TERRIFNSFL string
	TERRIFNSUL string
}

//CreateTable ...
func (a *AddressObject) CreateTable(db *pg.DB) {
	createTable(a, db)
}

//DropTable ...
func (a *AddressObject) DropTable(db *pg.DB) {
	dropTable(a, db)
}

//DropNotRelevantRows ...
func (a *AddressObject) DropNotRelevantRows(db *pg.DB) {
	_, err := db.Model(a).Where("livestatus != 1 AND currstatus != 0").Delete()
	if err != nil {
		logger.Error("Removed irrelevant rows",
			zap.String("table", "addrobj"),
			zap.Error(err),
		)
	}
	logger.Info("Removed irrelevant rows",
		zap.String("table", "addrobj"),
	)
}

//CreateIndexes ...
func (a *AddressObject) CreateIndexes(db *pg.DB) {
	_, err := db.Model(a).Exec(`			
			-- primary key
			-- ALTER TABLE fias.addrobj ADD CONSTRAINT addrobj_pkey PRIMARY KEY(aoid);

			--  create btree indexes
			CREATE INDEX aoguid_pk_idx ON fias.addrobj USING btree (aoguid);
			CREATE UNIQUE INDEX aoid_idx ON fias.addrobj USING btree (aoid);
			CREATE INDEX parentguid_idx ON fias.addrobj USING btree (parentguid);
			CREATE INDEX currstatus_idx ON fias.addrobj USING btree (currstatus);
			CREATE INDEX aolevel_idx ON fias.addrobj USING btree (aolevel);
			CREATE INDEX formalname_idx ON fias.addrobj USING btree (formalname);
			CREATE INDEX offname_idx ON fias.addrobj USING btree (offname);
			CREATE INDEX shortname_idx ON fias.addrobj USING btree (shortname);
			CREATE INDEX shortname_aolevel_idx ON fias.addrobj USING btree (shortname, aolevel);

			-- trigram indexes to speed up text searches
			CREATE INDEX formalname_trgm_idx on fias.addrobj USING gin (formalname gin_trgm_ops);
			CREATE INDEX offname_trgm_idx on fias.addrobj USING gin (offname gin_trgm_ops);
		`)
	/*
		-- foreign key
		ALTER TABLE fias.addrobj
		ADD CONSTRAINT addrobj_parentguid_fkey FOREIGN KEY (parentguid)
		REFERENCES fias.addrobj (aoguid) MATCH SIMPLE
		ON UPDATE CASCADE ON DELETE NO ACTION;
	*/

	if err != nil {
		logger.Error("Create indexes error: ",
			zap.String("table", "addrobj"),
			zap.Error(err),
		)
	}
	logger.Info("Index created",
		zap.String("table", "addrobj"),
	)
}

//Parse ...
func (a *AddressObject) Parse(folder string, db *pg.DB) {
	parseFolder(folder, a, db)
}
