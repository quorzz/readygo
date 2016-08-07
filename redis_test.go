package readygo

import (
	"testing"
)

var redis *Redis
var netErr error

func init() {

	config := &Config{
		Password:    "123456", // your passworf of redis, set "" as no password
		Database:    10,
		PoolMaxIdle: 10,
	}
	redis, netErr = Dial(config)

	redis.PingOnPool()
}

func TestCommand(t *testing.T) {

	msg, err := redis.Execute("set", "a", "123")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if ok, err := msg.Bool(); err != nil || !ok {
		t.FailNow()
	}

	msg, err = redis.Execute("get", "a")
	if err != nil {
		t.Error(err)
	}

	if str, err := msg.String(); err != nil || str != "123" {
		t.Log(msg)
		t.FailNow()
	}

	msg, err = redis.Execute("GET", "90modfsdgp975ksbksgk ")

	if err != nil {
		t.Error(err)
	}

	if str, err := msg.String(); err != nil {
		t.Error(err)
	} else if str != "" {
		t.Log(str)
		t.FailNow()
	}
}

func TestError(t *testing.T) {
	msg, err := redis.Execute("INVALID_COMMAND", "bbb009999")

	if err != nil {
		t.Error(err)
	} else if !msg.HasError() {
		t.Log("should be errors here")
		t.FailNow()
	}
}
