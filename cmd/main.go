package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lukisxyz/nexmed-service/internal/auth"
	mw "github.com/lukisxyz/nexmed-service/internal/middleware"
	"github.com/lukisxyz/nexmed-service/internal/profile"
	jwt "github.com/lukisxyz/nexmed-service/internal/utils/token"
	"github.com/lukisxyz/nexmed-service/lib/config"

	_ "github.com/lukisxyz/nexmed-service/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// @title User management
// @version 1.0
// @description Simple user management
// @host http://103.175.217.181 localhost:8080
// @schemes http https
// @BasePath /api
func main() {
	ctx := context.Background()
	
	// load configuration
	configFileName := ""
	flag.StringVar(&configFileName, "config", "config.yaml", "Config file name")
	flag.Parse()
	cfg := config.DefaultConfig()
	cfg.LoadFromEnv()
	if len(configFileName) > 0 {
		err := config.LoadConfigFromFile(configFileName, &cfg)
		if err != nil {
			log.Warn().Str("file", configFileName).Err(err).Msg("cannot load config file, use defaults")
		}
	}

	// create connection postgres
	pool, err := pgxpool.New(ctx, cfg.DBConfig.ConnStr())
	if (err != nil) {
		log.Error().Err(err).Msg("unable to connect to postgres database")
	}

	// create connection redis
	rdb := redis.NewClient(&redis.Options{
        Addr:     cfg.RedisConfig.Addr(),
        Password: cfg.RedisConfig.Password,
        DB:       0,
    })
	defer rdb.Close()

	auth.SetConnection(pool, rdb);
	profile.SetConnection(pool, rdb);

	// set jwt token
	jwt.SetJwtConfig(cfg.JWTConfig.Secret)

	// rate limiter
	mw.SetRateLimit(100, time.Minute)

	// set router using chi
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // Allow all headers
		ExposedHeaders:   []string{"*"}, // Expose all headers
		AllowCredentials: false,         // Must be false when using "*"
		MaxAge:           86400,         // Cache preflight response for 24 hours
	}))
	r.Use(mw.RateLimitMiddleware)
    r.Get("/swagger/*", httpSwagger.Handler(
        httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
    ))

	r.Mount("/api/", auth.Router())
	r.Mount("/api/profile", profile.Router())

	// generate server instance
	server := &http.Server{
		Addr:    cfg.Listener.Addr(),
		Handler: r,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal().Msg("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal().Err(err)
		}
		serverStopCtx()
	}()

	// Run the server
	fmt.Println("Server is running on port 8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

