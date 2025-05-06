package strategy

import (
	"RateBalancer/internal/service/balancer"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestRoundRobin_GetNext(t *testing.T) {
	tests := []struct {
		name               string
		aliveStates        []bool
		initCurrent        uint64
		expectedBackendURL string
		expectedErr        error
		expectedCurrent    uint64
	}{
		{
			name:               "all alive",
			aliveStates:        []bool{true, true, true},
			initCurrent:        1,
			expectedBackendURL: "http://backend2",
			expectedErr:        nil,
			expectedCurrent:    2,
		},
		{
			name:               "only second one alive",
			aliveStates:        []bool{false, true, false},
			initCurrent:        0,
			expectedBackendURL: "http://backend1",
			expectedErr:        nil,
			expectedCurrent:    1,
		},
		{
			name:               "check move idx if next is dead",
			aliveStates:        []bool{false, true, false},
			initCurrent:        1,
			expectedBackendURL: "http://backend1",
			expectedErr:        nil,
			expectedCurrent:    1,
		},
		{
			name:               "all service unavailable error",
			aliveStates:        []bool{false, false, false},
			initCurrent:        0,
			expectedBackendURL: "",
			expectedErr:        AllServiceNotAvailable,
			expectedCurrent:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backends := make([]*balancer.Backend, len(tt.aliveStates))
			for i, alive := range tt.aliveStates {
				u, err := url.Parse(fmt.Sprintf("http://backend%d", i))
				if err != nil {
					t.Fatalf("failed to parse URL: %v", err)
				}

				b := &balancer.Backend{URL: u}
				b.SetAlive(alive)
				backends[i] = b
			}
			backendPool := &balancer.BackendPool{
				Backends: backends,
			}
			backendPool.Current.Store(tt.initCurrent)

			rr := &RoundRobin{bp: backendPool}

			backend, err := rr.GetNext()

			assert.ErrorIs(t, err, tt.expectedErr)

			if tt.expectedBackendURL != "" {
				assert.Equal(t, tt.expectedBackendURL, backend.URL.String())
			} else {
				assert.Nil(t, nil, backend)
			}

			assert.Equal(t, tt.expectedCurrent, backendPool.Current.Load())
		})
	}
}
