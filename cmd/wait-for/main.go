package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/dmytro-boichenko/wait-for/internal/waiter"
)

const (
	defaultTimeout = time.Second
	defaultLimit   = 30

	maximumTimeoutInSeconds = 180
	maximumLimit            = 100
)

var (
	serviceFlagName = "service"
	timeoutFlagName = "timeout"
	limitFlagName   = "limit"
)

func main() {
	app := cli.App{
		Name:  "wait-for",
		Usage: "used for controlling in docker-compose builds for correct waiting for required resources like databases",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "service",
				Aliases:  []string{"s"},
				Usage:    fmt.Sprintf("service name. %s", waiter.NamesMessage()),
				Required: true,
			},
			&cli.IntFlag{
				Name:    "timeout",
				Aliases: []string{"t"},
				Usage:   "timeout in seconds between repeats",
				Value:   int(defaultTimeout.Seconds()),
			},
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"l"},
				Usage:   "number of repeats",
				Value:   defaultLimit,
			},
		},
		Action: func(c *cli.Context) error {
			service := c.String(serviceFlagName)

			start := time.Now()
			t := timeout(c)
			l := limit(c)

			ready, err := waiter.Wait(service, t, l)
			if err != nil {
				return err
			}

			message := resultMessage(service, start, ready)
			if ready {
				log.Print(message)
			} else {
				return fmt.Errorf(message)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func timeout(c *cli.Context) time.Duration {
	t := c.Int(timeoutFlagName)

	if t > 0 && t < maximumTimeoutInSeconds {
		return time.Duration(t) * time.Second
	}

	return defaultTimeout
}

func limit(c *cli.Context) int {
	l := c.Int(limitFlagName)

	if l > 0 && l < maximumLimit {
		return l
	}

	return defaultLimit
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
