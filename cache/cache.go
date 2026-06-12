package cache

import (
	"context"
	"time"
	"golang.org/x/sync/singleflight"
)

type Cache struct {
	sfGroup singleflight.Group
	Options Options
}

func (c *Cache) GetOrSet(ctx context.Context, key string, dest interface{}, ttl time.Duration, fetch func(ctx context.Context) (interface{}, error)) error {
	if !c.Options.EnableSingleflight {
		return c.executeFetch(ctx, key, dest, ttl, fetch)
	}

	val, err, _ := c.sfGroup.Do(key, func() (interface{}, error) {
		data, fetchErr := fetch(ctx)
		if fetchErr != nil {
			return nil, fetchErr
		}
		_ = c.Set(ctx, key, data, ttl)
		return data, nil
	})

	if err != nil {
		return err
	}

	return c.assignValue(dest, val)
}

func (c *Cache) executeFetch(ctx context.Context, key string, dest interface{}, ttl time.Duration, fetch func(ctx context.Context) (interface{}, error)) error {
	data, err := fetch(ctx)
	if err != nil {
		return err
	}
	_ = c.Set(ctx, key, data, ttl)
	return c.assignValue(dest, data)
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error { return nil }
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error { return nil }
func (c *Cache) assignValue(dest interface{}, val interface{}) error { return nil }
