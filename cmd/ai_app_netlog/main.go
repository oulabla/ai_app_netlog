package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/oulabla/ai_app_netlog/internal/app/repository"
	"github.com/oulabla/ai_app_netlog/internal/app/service"
	"github.com/oulabla/ai_app_netlog/internal/config"
	"github.com/oulabla/ai_app_netlog/internal/config/secret"
	"github.com/oulabla/ai_app_netlog/internal/metric"
	"github.com/oulabla/ai_app_netlog/internal/server"

	_ "github.com/oulabla/ai_app_netlog/internal/endpoints"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type appNameKeyType string

const appNameKey appNameKeyType = "application_name"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ────────────────────────────────────────────────
	// Флаги
	// ────────────────────────────────────────────────
	var (
		useLocalConfig bool
		debug          bool
	)

	flag.BoolVar(&useLocalConfig, "local", false, "use config/local.yaml instead of config/prod.yaml")
	flag.BoolVar(&debug, "debug", false, "enable debug logging and colored console output")
	flag.Parse()

	// ────────────────────────────────────────────────
	// Инициализация логгера (JSON в prod, цветной в debug)
	// ────────────────────────────────────────────────
	server.Init(debug)

	// Устанавливаем уровень лога глобально
	if debug {
		log.Logger = log.Logger.Level(zerolog.DebugLevel)
	}

	// Логируем запуск приложения
	log.Info().
		Bool("debug", debug).
		Str("config", func() string {
			if useLocalConfig {
				return "local.yaml"
			}
			return "prod.yaml"
		}()).
		Msg("starting application")

	// ────────────────────────────────────────────────
	// Конфигурация
	// ────────────────────────────────────────────────
	configFile := "config/prod.yaml"
	if useLocalConfig {
		configFile = "config/local.yaml"
	}

	configProvider, err := config.NewYAMLProvider(configFile)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("config_file", configFile).
			Msg("failed to load config")
	}

	config.SetProvider(configProvider)

	appName := config.GetString(ctx, config.K.ApplicationName)
	if appName == "" {
		log.Fatal().Msg("application name is not set")
	}
	ctx = context.WithValue(ctx, appNameKey, appName)

	// ────────────────────────────────────────────────
	// Секреты
	// ────────────────────────────────────────────────
	secretProvider, err := secret.NewYAMLSecretProvider("config/secret.yaml")
	if err != nil {
		log.Fatal().
			Err(err).
			Str("config_file", configFile).
			Msg("failed to load secrets")
	}
	secret.SetProvider(secretProvider)

	// ────────────────────────────────────────────────
	// Подключаем зависимости DI
	// ────────────────────────────────────────────────
	pgConn, err := CreatePostgres(ctx)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("pg connnection pool failed")
	}
	defer pgConn.Close()

	repo := repository.NewRepository(pgConn)
	appService := service.NewService(repo)
	server.GetInjector().Set(config.AppService, appService)

	// ────────────────────────────────────────────────
	// Prometheus metrics endpoint
	// ────────────────────────────────────────────────
	http.Handle("/metrics", promhttp.Handler())

	metricAddr := config.GetString(ctx, config.K.ServerMetricPort)

	go func() {
		log.Info().
			Str("addr", metricAddr).
			Msg("starting HTTP metrics server")

		if err := http.ListenAndServe(metricAddr, nil); err != nil {
			log.Error().
				Err(err).
				Str("addr", metricAddr).
				Msg("HTTP metrics server failed")
		}
	}()

	// ────────────────────────────────────────────────
	// gRPC сервер
	// ────────────────────────────────────────────────
	grpcAddr := config.GetString(ctx, config.K.ServerGrpcPort)

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("addr", grpcAddr).
			Msg("failed to listen on gRPC address")
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(metric.UnaryServerInterceptor()),
		grpc.UnaryInterceptor(server.UnaryErrorInterceptor),
	)

	server.RegisterAllGRPC(grpcServer)

	go func() {
		log.Info().
			Str("addr", grpcAddr).
			Msg("starting gRPC server")

		if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			log.Error().
				Err(err).
				Str("addr", grpcAddr).
				Msg("gRPC server failed")
		}
	}()

	// ────────────────────────────────────────────────
	// HTTP/JSON gateway (REST API)
	// ────────────────────────────────────────────────
	gwmux := runtime.NewServeMux()

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := server.RegisterAllGateway(ctx, gwmux, dialOpts); err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to register gateway endpoints")
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}).Handler(gwmux)

	httpAddr := config.GetString(ctx, config.K.ServerHttpPort)

	go func() {
		log.Info().
			Str("addr", httpAddr).
			Msg("starting HTTP/JSON gateway")

		if err := http.ListenAndServe(httpAddr, corsHandler); err != nil {
			log.Error().
				Err(err).
				Str("addr", httpAddr).
				Msg("HTTP/JSON gateway failed")
		}
	}()

	// ────────────────────────────────────────────────
	// Swagger UI (если реализован)
	// ────────────────────────────────────────────────
	go server.StartSwaggerServer(ctx)

	// ────────────────────────────────────────────────
	// Graceful shutdown
	// ────────────────────────────────────────────────
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	<-sigChan

	log.Info().Msg("received shutdown signal, gracefully stopping servers...")

	grpcServer.GracefulStop()
	cancel()

	log.Info().Msg("shutdown complete")
}

func CreatePostgresConnection(ctx context.Context) *pgx.Conn {
	pgURL := secret.Get(ctx, secret.PgCategory, secret.PgURL)
	pgConn, err := pgx.Connect(context.Background(), pgURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка подключения:")
	}

	// Проверка подключения
	var version string
	err = pgConn.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка запроса:")
	}
	log.Info().Msgf("PostgreSQL версия: %s", version)

	return pgConn
}

// Возвращает соединение и ошибку
func CreatePostgres(ctx context.Context) (*pgxpool.Pool, error) {
	pgURL := secret.Get(ctx, secret.PgCategory, secret.PgURL)

	config, err := pgxpool.ParseConfig(pgURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось разобрать строку подключения: %w", err)
	}

	// Важные настройки пула (подбери под нагрузку)
	config.MaxConns = 20 // максимум открытых соединений
	config.MinConns = 2  // минимум всегда готовых
	config.MaxConnLifetime = 45 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute // как часто проверять idle-соединения

	// можно добавить AfterConnect для логирования или инициализации
	// Правильная сигнатура AfterConnect
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		log.Info().Msg("Новое соединение в пуле создано")
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать пул: %w", err)
	}

	// Проверка подключения (опционально, но полезно при старте)
	var version string
	err = pool.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("тестовый запрос провалился: %w", err)
	}
	log.Info().Str("pg_version", version).Msg("PostgreSQL подключён успешно")

	return pool, nil
}
