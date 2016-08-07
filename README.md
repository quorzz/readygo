# readygo
a redis client for golang

Document
--------
- [GoDoc](http://godoc.org/github.com/quorzz/readgo)

Getting Started
------
```
    config := &redis.Config{
        Password:    "123456", // passworf of redis, set "" as no auth
        Database:    10,
        PoolMaxIdle: 10, // capacity of connection pool
    }
    redis, _ = redis.Dial(config)

    // PING when take connection from the pool
    redis.PingOnPool()

    message, _ := redis.Execute("SET", "abc", "demo")

    message.String()
```

more introduce about the message : [redis-protocol](http://github.com/quorzz/redis-protocol)