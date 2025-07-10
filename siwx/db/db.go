package db

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
)

var Connection *gorm.DB

func Connect(dsn string) {
    db, err := gorm.Open(postgres.New(postgres.Config{
        DSN:                  dsn,
        PreferSimpleProtocol: true,
    }), &gorm.Config{
        SkipDefaultTransaction: true,
        PrepareStmt:            false,
    })

    if err != nil {
        log.Fatal("Error connecting to the database:", err)
    }

    Connection = db
}
