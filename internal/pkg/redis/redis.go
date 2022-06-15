package redis

import (
	"context"
	"fmt"
	red "github.com/go-redis/redis/v8"
	"github.com/uibricks/studio-engine/internal/pkg/config"
	"strconv"
	"time"
)

type RedisClient interface {
	Set(key string, value string) (string, error)
	Get(key string) (string, error)
	Delete(key string) (string, error)
}

var host, port, pwd string
var db int

type Redis struct {
	client *red.Client
}

func ProvideDefaultRedisConn(config config.RedisConfig) (*Redis, error) {
	host = config.Host
	port = config.Port
	pwd = config.Password
	db, _ = strconv.Atoi(config.Db)
	return Conn(host, port, pwd, db)
}

func SetRedisClient(c *red.Client) *Redis {
	return &Redis{client: c}
}

func (r *Redis) DoString(ctx context.Context, action string, args ...interface{}) (string, error) {
	result, err := r.Do(ctx, action, args...)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result), nil
}

func (r *Redis) Do(ctx context.Context, action string, args ...interface{}) (interface{}, error) {
	if r == nil {
		var err error
		r, err = Conn(host, port, pwd, db)
		if err != nil {
			return "", err
		}
	}

	if err := validateConnection(r.client); err != nil {
		return "", &CreateDatabaseError{err: err}
	}

	cArgs := append([]interface{}{action}, args...)
	status := r.client.Do(context.Background(), cArgs...)
	if status.Err() != nil {
		return "", GenerateError(ctx, action, status.Err())
	}

	return status.Result()

}

func IsNil(err error) bool {
	return err == nil || err == red.Nil
}

func validateConnection(c *red.Client) error {
	_, err := c.Ping(context.Background()).Result() //makes sure database is connected
	if err != nil {
		return err
	}
	return nil
}

func Conn(host string, port string, pwd string, db int) (*Redis, error) {
	c := red.NewClient(&red.Options{
		Addr:       fmt.Sprintf("%s:%s", host, port),
		Password:   pwd,
		DB:         db,
		MaxRetries: 3,
	})

	if err := validateConnection(c); err != nil {
		return nil, &CreateDatabaseError{err: err}
	}

	return &Redis{client: c}, nil
}

func (r *Redis) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	status := r.client.Set(context.Background(), key, value, expiration)
	return GenerateError(ctx, "set", status.Err())
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	status := r.client.Get(context.Background(), key)
	if status.Err() != nil {
		return "", GenerateError(ctx, "get", status.Err())
	}
	value, _ := status.Result()
	return value, nil
}

func (r *Redis) HGet(ctx context.Context, client string, mode Mode, user string, key string) (string, error) {
	return r.DoString(ctx, "hget", PrepareKey(client, mode, user), key)
}

func (r *Redis) HGetAll(ctx context.Context, client string, mode Mode, user string) (interface{}, error) {
	return r.Do(ctx, "hgetall", PrepareKey(client, mode, user))
}

func (r *Redis) HSet(ctx context.Context, client string, mode Mode, user string, key string, value string) error {
	_, err := r.DoString(ctx, "hset", PrepareKey(client, mode, user), key, value)
	return err
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	status := r.client.Del(context.Background(), key)
	return GenerateError(ctx, "delete", status.Err())
}

func (r *Redis) HDel(ctx context.Context, client string, mode Mode, user string, key string) error {
	_, err := r.DoString(ctx, "hdel", PrepareKey(client, mode, user), key)
	return err
}