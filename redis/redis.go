package redis

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
)

var p *redis.Pool

// Config ...
type Config struct {
	IP          string        `yaml:"ip"`
	Port        string        `yaml:"port"`
	MaxIdle     int           `yaml:"maxIdle"`
	MaxActive   int           `yaml:"maxActive"`
	IdleTimeout time.Duration `yaml:"idleTimeout"`
	Wait        bool          `yaml:"wait"`
}

// Init ...
func Init(c *Config) error {
	p = &redis.Pool{
		MaxIdle:     c.MaxIdle,
		MaxActive:   c.MaxActive,
		IdleTimeout: c.IdleTimeout,
		Wait:        c.Wait,
		Dial: func() (conn redis.Conn, e error) {
			c, err := redis.Dial("tcp", c.IP+":"+c.Port)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	if _, err := p.Dial(); err != nil {
		return err
	}

	return nil
}

// GetPoolStats ...
func GetPoolStats() redis.PoolStats {
	return p.Stats()
}

// GetConn ...
func GetConn() redis.Conn {
	return p.Get()
}

// Exists ...
func Exists(k string) (bool, error) {
	c := p.Get()
	defer c.Close()

	if c.Err() != nil {
		return false, c.Err()
	}

	return redis.Bool(c.Do("EXISTS", k))
}

// SetExJSON ...
func SetExJSON(k string, v interface{}, ex int) error {
	c := p.Get()
	defer c.Close()
	if c.Err() != nil {
		return c.Err()
	}

	data, _ := json.Marshal(v)
	_, err := c.Do("SETEX", k, ex, data)
	return err
}

// GetJSON ...
func GetJSON(k string) ([]byte, error) {
	c := p.Get()
	defer c.Close()
	if c.Err() != nil {
		return nil, c.Err()
	}

	return redis.Bytes(c.Do("GET", k))
}

// SetExString ...
func SetExString(k, v string, ex int) error {
	c := p.Get()
	defer c.Close()
	if c.Err() != nil {
		return c.Err()
	}

	_, err := c.Do("SETEX", k, ex, v)
	return err
}

// GetString ...
func GetString(k string) (string, error) {
	c := p.Get()
	defer c.Close()
	if c.Err() != nil {
		return "", c.Err()
	}

	return redis.String(c.Do("GET", k))
}

// DelKey ...
func DelKey(k string) error {
	c := p.Get()
	defer c.Close()

	if c.Err() != nil {
		return c.Err()
	}

	_, err := c.Do("DEL", k)
	return err
}
