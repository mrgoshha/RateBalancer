package tests

import (
	"RateBalancer/internal/handler/http/balancer"
	"RateBalancer/internal/handler/http/limiter"
	serviceBalancer "RateBalancer/internal/service/balancer"
	serviceLimiter "RateBalancer/internal/service/limiter"
	"RateBalancer/pkg/hash"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"testing"
	"time"
)

type TestSuite struct {
	suite.Suite
	log      *slog.Logger
	db       *sqlx.DB
	servers  []*httptest.Server
	urls     []*url.URL
	balancer *balancer.Balancer
	limiter  *limiter.Limiter
}

func TestRateLimiterSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	dataSource := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		"test", "test", "localhost", "5432", "rateLimiterTest")

	if db, err := sqlx.Connect("postgres", dataSource); err != nil {
		s.FailNow("Failed to connect to postgres", err)
	} else {
		s.db = db
	}

	s.initDeps()

	if err := s.initDB(); err != nil {
		s.FailNow("Failed to create and populate DB", err)
	}
}

func (s *TestSuite) TearDownSuite() {
	for _, srv := range s.servers {
		srv.Close()
	}
	s.db.Close()
}

func (s *TestSuite) initDeps() {

	s.log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	n := 3
	s.servers = make([]*httptest.Server, n)
	s.urls = make([]*url.URL, n)
	for i := 0; i < n; i++ {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(fmt.Sprintf("Server%d", i)))
		}))
		s.servers[i] = srv

		u, err := url.Parse(srv.URL)
		if err != nil {
			s.FailNow("Failed to parse url", err)
		}
		s.urls[i] = u
	}

	ls := serviceLimiter.NewServiceLimiter(s.db, hash.NewSHA1Hasher(), &serviceLimiter.Config{Capacity: 100, PerSecond: 10})
	s.limiter = limiter.NewLimiter(ls, s.log)
}

func (s *TestSuite) newBackendPool(urls []*url.URL) *serviceBalancer.BackendPool {
	pool := &serviceBalancer.BackendPool{
		Backends: make([]*serviceBalancer.Backend, len(urls)),
	}
	pool.Current.Store(0)

	for i, u := range urls {
		proxy := httputil.NewSingleHostReverseProxy(u)

		backend := &serviceBalancer.Backend{
			URL:          u,
			ReverseProxy: proxy,
		}

		backend.SetAlive(true)

		pool.Backends[i] = backend
	}

	return pool
}

func (s *TestSuite) initDB() error {
	insertClient := `INSERT INTO clients VALUES ('2fa85ba2-2301-4ded-a456-ffd922271089',
	                    					'0c2975ce3a7af8aef09a40655c38129822c5d074',
	                    					0,
	                    					$1,
	                    					2,
	                    					1);`

	_, err := s.db.Exec(insertClient, time.Now().UTC())
	if err != nil {
		s.FailNow("Failed to insert client", err)
	}

	return nil
}
