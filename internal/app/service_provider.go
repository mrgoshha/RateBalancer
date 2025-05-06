package app

import (
	"RateBalancer/config"
	clientRepo "RateBalancer/internal/adapter/dbs/postgres/client"
	"RateBalancer/internal/handler/http/api"
	"RateBalancer/internal/handler/http/balancer"
	"RateBalancer/internal/service"
	balancerService "RateBalancer/internal/service/balancer"
	"RateBalancer/internal/service/balancer/strategy"
	clientService "RateBalancer/internal/service/client"
	"RateBalancer/internal/service/limiter"
	"RateBalancer/pkg/hash"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
)

type ServiceProvider struct {
	log              *slog.Logger
	config           *config.Config
	db               *sqlx.DB
	balancer         *balancer.Balancer
	strategy         service.Strategy
	backendPool      *balancerService.BackendPool
	router           *http.ServeMux
	adminRouter      *http.ServeMux
	limiterService   *limiter.Service
	clientRepository *clientRepo.Repository
	clientService    *clientService.Service
	clientController *api.ClientController
	hasher           *hash.SHA1Hasher
}

func NewServiceProvider(log *slog.Logger, db *sqlx.DB, cfg *config.Config) *ServiceProvider {
	return &ServiceProvider{
		log:    log,
		db:     db,
		config: cfg,
	}
}

func (s *ServiceProvider) HttpRouter() *http.ServeMux {
	if s.router == nil {
		s.router = http.NewServeMux()
	}
	return s.router
}

func (s *ServiceProvider) HttpAdminRouter() *http.ServeMux {
	if s.adminRouter == nil {
		s.adminRouter = http.NewServeMux()
	}
	return s.adminRouter
}

func (s *ServiceProvider) Strategy() (service.Strategy, error) {
	if s.strategy == nil {
		switch s.config.LoadBalancer.Strategy {
		case service.RoundRobin:
			s.strategy = strategy.NewRoundRobinBalancer(s.BackendPool())
		case service.Random:
			s.strategy = strategy.NewRandomBalancer(s.BackendPool())
		default:
			return nil, fmt.Errorf("unknown strategy: %s", s.config.LoadBalancer.Strategy)
		}
	}

	return s.strategy, nil
}

func (s *ServiceProvider) BackendPool() *balancerService.BackendPool {
	if s.backendPool == nil {
		s.backendPool = balancerService.NewBackendPool(s.config.LoadBalancer)
	}
	return s.backendPool
}

func (s *ServiceProvider) LimiterService() *limiter.Service {
	if s.limiterService == nil {
		s.limiterService = limiter.NewServiceLimiter(s.db, s.Hash(), s.config.Limiter)
	}
	return s.limiterService
}

func (s *ServiceProvider) ClientRepository() *clientRepo.Repository {
	if s.clientRepository == nil {
		s.clientRepository = clientRepo.NewRepository(s.db)
	}
	return s.clientRepository
}

func (s *ServiceProvider) Hash() *hash.SHA1Hasher {
	if s.hasher == nil {
		s.hasher = hash.NewSHA1Hasher()
	}
	return s.hasher
}

func (s *ServiceProvider) ClientService() *clientService.Service {
	if s.clientService == nil {
		s.clientService = clientService.NewService(s.ClientRepository(), s.Hash(), s.config.Limiter)
	}
	return s.clientService
}

func (s *ServiceProvider) ClientController() *api.ClientController {
	if s.clientController == nil {
		s.clientController = api.NewClientController(s.log, s.ClientService(), s.HttpAdminRouter())
	}
	return s.clientController
}

func (s *ServiceProvider) RegisterControllers() {
	s.ClientController()
}
