package main

import (
	"log"
	"os"

	"github.com/garyburd/redigo/redis"
	minio "github.com/minio/minio/cmd"

	// Import gateway
	_ "github.com/minio/minio/cmd/gateway"
)

func main() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatalln("REDIS_URL cannot be empty")
	}

	connPool := redis.NewPool(func() (redis.Conn, error) {
		return redis.DialURL(redisURL)
	}, 3)

	c := connPool.Get()
	if _, err := c.Do("PING"); err != nil {
		c.Close()
		log.Fatalf("failed to connect to redis, %v", err)
	}
	c.Close()

	minio.GlobalCredentialProvider = &credSyncer{connPool}

	if len(os.Args) == 1 {
		// start from env config
		mode := os.Getenv("MINIO_MODE")
		switch mode {
		case "gateway":
			gateway := os.Getenv("MINIO_GATEWAY")
			endpoint := os.Getenv("MINIO_GATEWAY_ENDPOINT")
			os.Args = append(os.Args, "gateway", gateway, endpoint)
		case "server":
			os.Args = append(os.Args, os.Getenv("MINIO_EXPORT"))
		default:
			log.Fatal(`len(os.Args)==1, but valid mode "gateway" or "server" is not provided`)
		}
	}

	// start
	minio.Main(os.Args)
}
