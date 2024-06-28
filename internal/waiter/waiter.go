package waiter

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var (
	waiterConstructors = map[string]constructor{
		"mysql":         mySQLWaiter,
		"elasticsearch": elasticsearchWaiter,
		"redis":         redisWaiter,
		"postgres":      postgresWaiter,
	}
)

type constructor func() waiter

type waiter interface {
	waitFor() (bool, error)
	name() string
}

func Wait(serviceName string, timeout time.Duration, limit int) (bool, error) {
	c, ok := waiterConstructors[serviceName]
	if !ok {
		return false, fmt.Errorf("waiting for %s doesn't supported. %s", serviceName, NamesMessage())
	}

	w := c()
	for i := 0; i < limit; i++ {
		ready, err := w.waitFor()
		if err != nil {
			log.Println(err)
		}

		if ready {
			return true, nil
		}

		time.Sleep(timeout)
	}

	return false, nil
}

func envVar(name, defaultValue string) string {
	v := os.Getenv(name)
	if v == "" {
		return defaultValue
	}
	return v
}

func NamesMessage() string {
	values := make([]string, 0)
	for k := range waiterConstructors {
		values = append(values, k)
	}
	return fmt.Sprintf("Supported: %v", strings.Join(values, ", "))
}
