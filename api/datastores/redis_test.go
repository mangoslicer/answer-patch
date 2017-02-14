package datastores

import (
	"testing"
)

func TestConnectToRedis(t *testing.T) {
	if _, err := ConnectToRedis().Do("PING"); err != nil {
		t.Error(err)
	}
}
