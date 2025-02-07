package cache

import (
	"context"
	"fmt"
	"time"
)

func (c *cache) SetCompanyID(ctx context.Context, name, id string) (err error) {
	if cmd := c.redisClient.Set(ctx,
		getCompanyKey(name), id,
		time.Duration(time.Hour*24)); cmd.Err() != nil {
		return cmd.Err()
	}

	return
}

func (c *cache) GetCompanyID(ctx context.Context, name string) (string, error) {
	var (
		cmd = c.redisClient.Get(ctx, getCompanyKey(name))
	)

	if cmd.Err() != nil {
		return "", cmd.Err()
	}

	out, err := cmd.Result()
	if err != nil {
		return "", err
	}

	return out, nil
}

// func (c *cache) DeleteHero(ctx context.Context, key string) (err error) {
// 	cmd := c.redisClient.Del(ctx, getHeroKey(key))
// 	if cmd.Err() != nil {
// 		return cmd.Err()
// 	}
// 	return
// }

func getCompanyKey(name string) string {
	return fmt.Sprintf("company:%s", name)
}
