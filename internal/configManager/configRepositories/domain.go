package configrepositories

import (
	"context"
	"database/sql"
)

type ConfigDbConnection struct {
	*sql.DB
}

type SQLiteConfigRepository struct {
	ctx          context.Context
	dbConnection *ConfigDbConnection
}

func NewSQLiteConfigRepository(ctx context.Context, conn *ConfigDbConnection) *SQLiteConfigRepository {
	return &SQLiteConfigRepository{
		ctx:          ctx,
		dbConnection: conn,
	}
}
