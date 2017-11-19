package redis

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// DB redis DB connection
type DB struct {
	conn redis.Conn
}

// redisDb global db instance
var redisDb *DB

// Connect - makes a redis db connection
func Connect(url string) (*DB, error) {
	if redisDb == nil {
		redisDb = new(DB)
		var err error

		redisDb.conn, err = redis.DialURL(url, redis.DialConnectTimeout(30*time.Second))
		if err != nil {
			return nil, err
		}
	}

	return redisDb, nil
}

// SetValue sets the given key/value to db
func (redisDb *DB) SetValue(key string, value interface{}) error {
	_, err := redisDb.conn.Do("SET", key, value)
	return err
}

// GetValue return the value for a given key
func (redisDb *DB) GetValue(key string) (interface{}, error) {
	return redisDb.conn.Do("GET", key)
}

// SetStringValue sets the given key/value to db
func (redisDb *DB) SetStringValue(key string, value string, expiration ...interface{}) error {
	_, err := redisDb.conn.Do("SET", key, value)

	if err == nil && expiration != nil {
		redisDb.conn.Do("EXPIRE", key, expiration[0])
	}

	return err
}

// GetStringValue return the value for a given key
func (redisDb *DB) GetStringValue(key string) (string, error) {
	v, err := redisDb.conn.Do("GET", key)
	if err != nil {
		return "", nil
	}
	if v == nil {
		return "", nil
	}
	return redis.String(v, nil)
}
