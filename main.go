package main

import (
	"context"
	"flag"
	"net/http"
	"nextmed-service/account"
	"nextmed-service/auth"
	custom_middleware "nextmed-service/middleware"

	_ "github.com/lukisxyz/nexmed-service/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title User management
// @version 1.0
// @description Simple user management
// @host localhost:3000
// @BasePath /api
func main() {
	ctx := context.Background()
	
	// load configuration
	var configFileName string
	flag.StringVar(&configFileName, "config", "config.yaml", "Config file name")
	flag.Parse()
	cfg := loadConfig(configFileName)

	// create connection postgres
	pool, err := pgxpool.New(ctx, cfg.DBConfig.ConnStr())
	if (err != nil) {
		log.Error().Err(err).Msg("unable to connect to postgres database")
	}
	defer pool.Close()

	// create connection redis
	rdb := redis.NewClient(&redis.Options{
        Addr:     cfg.RedisConfig.Addr(),
        Password: cfg.RedisConfig.Password,
        DB:       0,
    })
	defer rdb.Close()

	if err := custom_middleware.SetConnectionRedis(rdb); err != nil {
		log.Error().Err(err).Msg("unable to connect to redis database")
	}

	// set router using chi
	r := chi.NewRouter()
    r.Get("/swagger/*", httpSwagger.Handler(
        httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
    ))
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Mount("/api/", auth.Router())
	r.Mount("/api/account", account.Router())

	log.Info().Msg("starting up server...")
	if err := http.ListenAndServe(cfg.Listener.Addr(), r); err != nil {
		log.Fatal().Err(err).Msg("Failed to start the server")
		return 
	}
	log.Info().Msg("server stopped")
}
