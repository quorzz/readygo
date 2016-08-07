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
	PoolChecker func(item interface{}) error
}

var DefaultConfig = &Config{
	Network:     "tcp",
	Address:     ":6379",
	Password:    "",
	Database:    0,
	Timeout:     6 * time.Second,
	PoolMaxIdle: 0,
	PoolChecker: func(item interface{}) error {
		return nil
	},
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

	if nil == config.PoolChecker {
		config.PoolChecker = DefaultConfig.PoolChecker
	}
}

func (r *Redis) newConnection() (*connection, error) {
	netConn, err := net.DialTimeout(r.network, r.address, r.timeout)
	if err != nil {
		return nil, err
	}
	conn := NewConnection(netConn)

	//  AUTH PASSWORD
	if "" != r.password {
		msg, err := conn.Execute("AUTH", r.password)
		if err != nil {
			return nil, err
		}
		if msg.HasError() {
			return nil, msg.Error
		}
	}

	// SELECT DATABASE
	if r.database > 0 {
		msg, err := conn.Execute("SELECT", r.database)
		if err != nil {
			return nil, err
		}
		if msg.HasError() {
			return nil, msg.Error
		}
	}
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

func (redis *Redis) PingOnPool() {
	redis.pool.Checker = func(item interface{}) error {
		conn, _ := item.(*connection)
		_, err := conn.Execute("PING")
		return err
	}
}

func (redis *Redis) Execute(args ...interface{}) (*protocol.Message, error) {

	item, err := redis.pool.Get()
	conn, _ := item.(*connection)

	defer func() { redis.pool.Put(conn) }()

	if err != nil {
		return nil, err
	}
	return conn.Execute(args...)
}
