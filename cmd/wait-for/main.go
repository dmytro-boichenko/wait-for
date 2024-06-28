package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dmytro-boichenko/wait-for/internal/waiter"
)

var (
	serviceFlag = flag.String("service", "", fmt.Sprintf("service name. Mandatory. %s", waiter.NamesMessage()))
	timeoutFlag = flag.Int("timeout", int(waiter.DefaultTimeout.Seconds()), "timeout in seconds between repeats")
	limitFlag   = flag.Int("limit", waiter.DefaultLimit, "number of repeats")
)

func main() {
	flag.Parse()

	start := time.Now()

	ready, err := waiter.Wait(*serviceFlag, *timeoutFlag, *limitFlag)
	if err != nil {
		log.Fatal(err)
	}

	message := resultMessage(*serviceFlag, start, ready)
	if ready {
		log.Print(message)
	} else {
		log.Fatal(message)
	}
}

func resultMessage(serviceName string, start time.Time, ready bool) string {
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
