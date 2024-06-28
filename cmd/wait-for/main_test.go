package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPrepareResultMessage(t *testing.T) {
	testCases := []struct {
		name   string
		start  time.Time
		ready  bool
		result string
	}{
		{
			name:   "ready_0s",
			start:  time.Now(),
			ready:  true,
			result: "Service test_waiter is ready",
		},
		{
			name:   "not_ready_0s",
			start:  time.Now(),
			ready:  false,
			result: "Service test_waiter is not ready",
		},
		{
			name:   "ready_17s",
			start:  time.Now().Add(-17 * time.Second),
			ready:  true,
			result: "Service test_waiter is ready in 17s",
		},
		{
			name:   "not_ready_17s",
			start:  time.Now().Add(-17 * time.Second),
			ready:  false,
			result: "Service test_waiter is not ready in 17s",
		},
		{
			name:   "ready_1h_23m_45secs",
			start:  time.Now().Add(-5025 * time.Second),
			ready:  true,
			result: "Service test_waiter is ready in 1h23m45s",
		},
		{
			name:   "not_ready_1h_23m_45secs",
			start:  time.Now().Add(-5025 * time.Second),
			ready:  false,
			result: "Service test_waiter is not ready in 1h23m45s",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := resultMessage("test_waiter", tc.start, tc.ready)
			require.Equal(t, tc.result, actual)
		})
	}
}
