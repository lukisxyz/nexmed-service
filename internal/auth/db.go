package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)


var (
	pool *pgxpool.Pool
	rdb *redis.Client
)

func SetConnection(_pool *pgxpool.Pool, _rdb *redis.Client) {
	if _pool == nil {
		log.Fatal().Msg("unable set null to pool")
	}
	pool = _pool
	if _rdb == nil {
		log.Fatal().Msg("unable set null to redis")
	}
	rdb = _rdb
}