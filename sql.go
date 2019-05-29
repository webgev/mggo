package mggo

import (
    "github.com/go-pg/pg"
    "github.com/go-pg/pg/orm"
)

var db *pg.DB

type dbLogger struct{}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
    s, _ := q.FormattedQuery()
    LogInfo("Вызов sql:\n", s)
}

// SQLOpen - connect sql
func SQLOpen() {
    if db != nil {
        return
    }
    cs := config.Section("database")
    db = pg.Connect(&pg.Options{
        User:     cs.Key("user").String(),
        Password: cs.Key("password").String(),
        Database: cs.Key("database").String(),
        Addr:     cs.Key("address").String(),
        Network:     cs.Key("network").String(),
    })
    db.AddQueryHook(dbLogger{})
    LogInfo("Connect Database")
}

// SQLClose - disconnect sql
func SQLClose() {
    if db == nil {
        return
    }
    db.Close()
    db = nil
    LogInfo("Disconnect Database")
}

// SQL is open sql and return pg.DB
func SQL() *pg.DB {
    if db == nil {
        SQLOpen()
    }
    return db
}

// Scan is pg.Scan
func Scan(i ...interface{}) orm.ColumnScanner {
    return pg.Scan(i...)
}

// CreateTable is create table
func CreateTable(models []interface{}) {
    for _, model := range models {
        err := db.CreateTable(model, &orm.CreateTableOptions{
            IfNotExists: true,
        })
        if err != nil {
            LogInfo(err)
        }
    }
}
