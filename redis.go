package readygo

import (
	"github.com/quorzz/pool"
	"github.com/quorzz/redis-protocol"
	"net"
	"time"
)

type Redis struct {
	network  string
	address  string
	database int
	password string
	timeout  time.Duration
	pool     *pool.Pool
}

type Config struct {
	Network     string
	Address     string
	Password    string
	Database    int
	Timeout     time.Duration
	PoolMaxIdle int // 连接池容量
}

var DefaultConfig = &Config{
	Network:     "tcp",
	Address:     ":6379",
	Password:    "",
	Database:    0,
	Timeout:     6 * time.Second,
	PoolMaxIdle: 0,
}

func parseConfig(config *Config) {
	if config.Network == "" {
		config.Network = DefaultConfig.Network
	}

	if config.Address == "" {
		config.Address = DefaultConfig.Address
	}

	if config.Timeout == 0 {
		config.Timeout = DefaultConfig.Timeout
	}

	if config.PoolMaxIdle == 0 {
		config.PoolMaxIdle = DefaultConfig.PoolMaxIdle
	}
}

func (r *Redis) newConnection() (*connection, error) {
	netConn, err := net.DialTimeout(r.network, r.address, r.timeout)
	if err != nil {
		return nil, err
	}
	conn := NewConnection(netConn)
	//  "AUTH"

	// "SELECT DATABASE"

	return conn, nil
}

func Dial(config *Config) (*Redis, error) {

	parseConfig(config)

	redis := &Redis{
		network:  config.Network,
		address:  config.Address,
		database: config.Database,
		password: config.Password,
		timeout:  config.Timeout,
	}

	redis.pool = pool.NewPool(func() (interface{}, error) {
		return redis.newConnection()
	}, config.PoolMaxIdle)

	return redis, nil
}

func (redis *Redis) Execute(args ...interface{}) (*protocol.Message, error) {

	item, err := redis.pool.Get()
	conn, _ := item.(*connection)

	defer func() { redis.pool.Put(conn) }()

	if err != nil {
		return nil, err
	}

	if err := conn.Send(args...); err != nil {
		return nil, err
	}

	if err := conn.Flush(); err != nil {
		return nil, err
	}

	if message, err := conn.Receive(); err != nil {
		return nil, err
	} else {
		return message, nil
	}
}
