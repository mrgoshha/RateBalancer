package app

import (
	"RateBalancer/config"
	postgres "RateBalancer/internal/adapter/dbs/postgres"
	"RateBalancer/internal/handler/http/adminserver"
	"RateBalancer/internal/handler/http/balancer"
	"RateBalancer/internal/handler/http/limiter"
	"RateBalancer/internal/handler/http/server"
	"RateBalancer/internal/service/balancer/healthchecker"
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	log             *slog.Logger
	server          *http.Server
	adminServer     *http.Server
	db              *sqlx.DB
	serviceProvider *ServiceProvider
	balancer        *balancer.Balancer
	limiter         *limiter.Limiter
	healthChecker   *healthchecker.HealthChecker
	config          *config.Config
	configPath      string
}

func NewApp(configPath string) (*App, error) {
	a := &App{
		configPath: configPath,
	}

	err := a.initDeps()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go a.healthChecker.HealthCheck(ctx)

	// Graceful Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go a.runServer()
	a.log.Info("server started")

	go a.runAdminServer()
	a.log.Info("server admin started")
	
	<-ctx.Done()
	a.log.Info("shutting down...")

	cancel()
	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Error("failed to stop server", slog.String("error", err.Error()))
	}

	if err := a.adminServer.Shutdown(ctx); err != nil {
		a.log.Error("failed to stop admin server", slog.String("error", err.Error()))
	}

	a.log.Info("servers stopped")

	if err := a.db.Close(); err != nil {
		a.log.Error("failed to stop storage", slog.String("error", err.Error()))
	}
}

func (a *App) initDeps() error {
	inits := []func() error{
		a.initConfig,
		a.initLogger,
		a.initDb,
		a.initServiceProvider,
		a.initBalancer,
		a.initLimiter,
		a.initHealthChecker,
		a.initHttpServer,
		a.initHttpAdminServer,
	}

	for _, f := range inits {
		err := f()
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig() error {
	cfg, err := config.Load(a.configPath)
	if err != nil {
		return err
	}
	a.config = cfg
	return nil
}

func (a *App) initLogger() error {
	a.log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return nil
}

func (a *App) initDb() error {
	db, err := postgres.New(a.config.Database)
	if err != nil {
		a.log.Error("failed to init storage", slog.String("error", err.Error()))
		return err
	}

	a.db = db

	return nil
}

func (a *App) initServiceProvider() error {
	a.serviceProvider = NewServiceProvider(a.log, a.db, a.config)
	return nil
}

func (a *App) initBalancer() error {
	strategy, err := a.serviceProvider.Strategy()
	if err != nil {
		return err
	}
	a.balancer = balancer.NewBalancer(strategy, a.serviceProvider.BackendPool(), a.log)

	return nil
}

func (a *App) initLimiter() error {
	a.limiter = limiter.NewLimiter(a.serviceProvider.LimiterService(), a.log)
	return nil
}

func (a *App) initHealthChecker() error {
	a.healthChecker = healthchecker.NewHealthChecker(a.serviceProvider.BackendPool(), a.config.HealthChecker)
	return nil
}

func (a *App) initHttpServer() error {
	a.balancer.RegisterBalancer(a.serviceProvider.HttpRouter())
	handlerWithLimiter := a.limiter.RegisterLimiter(a.serviceProvider.HttpRouter())
	srv := server.NewServer(a.config.Server, a.log, handlerWithLimiter)

	a.server = srv
	return nil
}

func (a *App) initHttpAdminServer() error {
	a.serviceProvider.RegisterControllers()
	srv := adminserver.NewServer(a.config.AdminServer, a.log, a.serviceProvider.HttpAdminRouter())

	a.adminServer = srv
	return nil
}

func (a *App) runServer() {
	if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		a.log.Error("failed to start server", slog.String("error", err.Error()))
	}
}

func (a *App) runAdminServer() {
	if err := a.adminServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		a.log.Error("failed to start admin server", slog.String("error", err.Error()))
	}
}
