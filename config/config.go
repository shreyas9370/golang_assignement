package config

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis/v8"
)

var (
	DB          *sql.DB
	RedisClient *redis.Client
	Ctx         context.Context
)
