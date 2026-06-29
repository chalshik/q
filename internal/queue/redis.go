package queue

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr string) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &RedisClient{client: client}, nil
}

func (r *RedisClient) NewQueue(key string) *RedisQueue {
	return &RedisQueue{client: r.client, key: key}
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

type RedisQueue struct {
	client *redis.Client
	key    string
}

func (q *RedisQueue) Enqueue(jobID string) error {
	return q.client.LPush(context.Background(), q.key, jobID).Err()
}

// Dequeue blocks until a job ID is available and returns it.
func (q *RedisQueue) Dequeue() string {
	for {
		result, err := q.client.BRPop(context.Background(), 5*time.Second, q.key).Result()
		if err == nil {
			return result[1]
		}
	}
}
