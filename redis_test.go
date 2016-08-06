package readygo

import (
	"testing"
)

func TestConnection(t *testing.T) {

	redis, netErr := Dial(DefaultConfig)

	if netErr != nil {
		t.Error(netErr)
	}

	redis.Execute("set", "a", "123")
	msg, err := redis.Execute("get", "a")
	if err != nil {
		t.Error(err)
	}

	if str, err := msg.String(); err != nil || str != "123" {
		t.Log(msg)
		t.Fail()
	}
}
