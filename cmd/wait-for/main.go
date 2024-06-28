package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	maximumTimeoutInSeconds = 180
	maximumLimit            = 100

	defaultTimeout = time.Second
	defaultLimit   = 30

	unknownServiceName = ""
)

type waiting struct {
	constructor waiterConstructor
	params      waitingParams
}

type waiterConstructor func() waiter

type waiter interface {
	waitFor() (bool, error)
	name() string
}

type waitingParams struct {
	timeout time.Duration
	limit   int
}

var (
	defaultWaitingParams = waitingParams{
		timeout: defaultTimeout,
		limit:   defaultLimit,
	}

	waiters = map[string]waiting{
		"mysql": {
			constructor: getMySQLWaiter,
			params:      defaultWaitingParams,
		},
		"elasticsearch": {
			constructor: getElasticsearchWaiter,
			params:      defaultWaitingParams,
		},
		"redis": {
			constructor: getRedisWaiter,
			params:      defaultWaitingParams,
		},
		"postgres": {
			constructor: getPostgresWaiter,
			params:      defaultWaitingParams,
		},
	}

	serviceFlag = flag.String("service", unknownServiceName, fmt.Sprintf("service name. Mandatory. %s", waiterNamesMessage()))
	timeoutFlag = flag.Int("timeout", int(defaultWaitingParams.timeout.Seconds()), "timeout in seconds between repeats")
	limitFlag   = flag.Int("limit", defaultWaitingParams.limit, "number of repeats")
)

func main() {
	flag.Parse()

	w, err := getWaiterFactory(*serviceFlag)
	if err != nil {
		log.Fatal(err)
	}

	waiter := w.constructor()
	params := parseWaiterParams(w.params)

	start := time.Now()
	ready := wait(waiter, params)

	resultMessage := prepareResultMessage(waiter.name(), start, ready)
	if ready {
		log.Print(resultMessage)
	} else {
		log.Fatal(resultMessage)
	}
}

func wait(w waiter, params waitingParams) bool {
	for i := 0; i < params.limit; i++ {
		ready, err := w.waitFor()
		if err != nil {
			log.Println(err)
		}

		if ready {
			return true
		}

		time.Sleep(params.timeout)
	}

	return false
}

func getWaiterFactory(s string) (waiting, error) {
	if s == unknownServiceName {
		flag.Usage()
		os.Exit(-1)
	}

	w, ok := waiters[s]
	if !ok {
		return waiting{}, errors.Errorf("Waiting for %s doesn't supported. %s", s, waiterNamesMessage())
	}

	return w, nil
}

func waiterNamesMessage() string {
	values := make([]string, 0)
	for k := range waiters {
		values = append(values, k)
	}
	return fmt.Sprintf("Possible values are: %v", values)
}

func parseWaiterParams(params waitingParams) waitingParams {
	timeout := params.timeout
	limit := params.limit

	if t := *timeoutFlag; t > 0 && t < maximumTimeoutInSeconds {
		timeout = time.Duration(t) * time.Second
	}

	if l := *limitFlag; l > 0 && l < maximumLimit {
		limit = l
	}

	return waitingParams{
		timeout: timeout,
		limit:   limit,
	}
}

func envVar(name, defaultValue string) string {
	v := os.Getenv(name)
	if v == "" {
		return defaultValue
	}
	return v
}

func prepareResultMessage(serviceName string, start time.Time, ready bool) string {
	duration := time.Since(start).Truncate(time.Second)
	suffix := ""
	if duration >= time.Second {
		suffix = fmt.Sprintf("in %s", duration.String())
	}

	var msgTemplate string
	if ready {
		msgTemplate = "Service %s is ready %s"
	} else {
		msgTemplate = "Service %s is not ready %s"
	}

	msg := fmt.Sprintf(msgTemplate, serviceName, suffix)
	return strings.TrimSpace(msg)
}
