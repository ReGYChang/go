package redis

import (
	"context"
	"fmt"
	"log"
	"nexdata/pkg/config"
	"testing"

	"github.com/go-redis/redis/v9"
)

func TestGetClient(t *testing.T) {
	ctx := context.Background()
	config.Redis.Host = "localhost:6379"
	config.Redis.DB = 0

	rdb, err := GetClient()
	if err != nil {
		t.Fatalf("get redis client error: %v", err)
	}

	err = rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		t.Fatalf("set redis key failed: %v", err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		t.Fatalf("get redis key failed: %v", err)
	}
	log.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		log.Println("key2 does not exist")
	} else if err != nil {
		log.Fatalf("get redis key2 failed: %v", err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}
