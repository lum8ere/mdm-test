package db_manager

import (
	"errors"
	"mdm/libs/4_common/smart_context"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbManager struct {
	db *gorm.DB
}

func NewDbManager(sctx smart_context.ISmartContext) (*DbManager, error) {
	databaseUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return nil, errors.New("DATABASE_URL is not set")
	}
	sctx.Debugf("DATABASE_URL: %s", databaseUrl)

	if sctx == nil {
		sctx = smart_context.NewSmartContext()
	}

	db, err := gorm.Open(
		postgres.Open(databaseUrl),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	result := &DbManager{
		db: db,
	}

	return result, nil
}

func (dbmanager *DbManager) GetGORM() *gorm.DB {
	return dbmanager.db.Session(&gorm.Session{NewDB: true})
}
