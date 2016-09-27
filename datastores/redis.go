package datastores

import (
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/patelndipen/AP1/settings"
)

func ConnectToRedis() redis.Conn {

	dsn := settings.GetRedisDSN()

	conn, err := redis.Dial("tcp", dsn.Addr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Do("AUTH", dsn.Password)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = conn.Do("PING"); err != nil {
		conn.Close()
		log.Fatal(err)
	}

	return conn
}
