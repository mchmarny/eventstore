package queue

import (
	"log"
	"github.com/mchmarny/myevents/pkg/utils"
	"github.com/adjust/rmq"
	"gopkg.in/redis.v3"
)

const (

	defaultRedisHost  = "redis.default.svc.cluster.local:6379"
	defaultRedisQueue  = "myevents"
)

var (
	redisHost = utils.MustGetEnv("REDIS_HOST", defaultRedisHost)
	redisPass = utils.MustGetEnv("REDIS_PASS", "")
	redisQueue = utils.MustGetEnv("REDIS_QUEUE", defaultRedisQueue)
	q rmq.Queue
)


// GetQueue initializes the redis queue
func GetQueue() rmq.Queue {

	// if queue already initialized
	if q != nil {
		return q
	}

	// redis client to set password
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPass,
		DB:       0,
	})

	// test
	log.Printf("Connecting to %s...", redisHost)
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Error on PING: %v", err)
	}

	// queue connection
	cnn := rmq.OpenConnectionWithRedisClient(defaultRedisQueue, client)

	// set queue
	q = cnn.OpenQueue(redisQueue)

	return q

}