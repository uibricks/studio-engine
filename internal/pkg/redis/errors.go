package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

// OperationError when cannot perform a given operation on database (SET,GET or DELETE)
type OperationError struct {
	operation string
	err       error
}

func (e *OperationError) Error() string {
	return fmt.Sprintf("Could not perform the %s operation. Error - %v", e.operation, e.err)
}

// CreateDatabaseError when cannot perform set on database
type CreateDatabaseError struct {
	err error
}

func (c *CreateDatabaseError) Error() string {
	return fmt.Sprintf("Could not create database. Error - %v", c.err)
}

// DownError when its not a redis.Nil response, in this case the database is down
type DownError struct {
	err error
}

func (d *DownError) Error() string {
	return fmt.Sprintf("Redis Database is down - %v", d.err)
}

func GenerateError(ctx context.Context, operation string, err error) error {
	if err != nil {
		if err == redis.Nil {
			return err
		}
		return &OperationError{operation: operation, err: err}
	}
	return nil
}
