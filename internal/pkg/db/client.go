package db

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/uibricks/studio-engine/internal/pkg/logger"
)

type DbClient struct{
	Connection *pg.DB
}

func ProvideDBClient() DbClient {
	return DbClient{}
}

func (client *DbClient) Connect(dbUrl,dbSchema string) error{
	options, err := pg.ParseURL(dbUrl)
	if err!=nil {
		logger.Sugar.Panicf("failed to parse the db url - %v", err)
		return err
	}

	options.OnConnect = func(ctx context.Context, conn *pg.Conn) error {
		_, err := conn.Exec("set search_path=?", dbSchema)
		if err != nil {
			logger.Sugar.Panicf("failed to set the default schema - %s", dbSchema)
			return err
		}
		return nil
	}

	options.MaxRetries = 3

	client.Connection = pg.Connect(options)
	client.Connection.AddQueryHook(dbLogger{})
	return nil
}

type dbLogger struct {
}

func (d dbLogger) BeforeQuery(ctx context.Context, q *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (d dbLogger) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	bytes, _ := q.FormattedQuery()
	fmt.Println(string(bytes))
	return nil
}

func (client *DbClient) Close() {
	if client.Connection!=nil{
		client.Connection.Close()
	}
}
