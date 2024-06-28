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

var (
	serviceFlagName = "service"
	timeoutFlagName = "timeout"
	limitFlagName   = "limit"
)

func main() {
	app := cli.App{
		Name:  "wait-for",
		Usage: "The wait-for tool is used for controlling in docker-compose builds for correct waiting for required resources like databases.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "service",
				Aliases:  []string{"s"},
				Usage:    fmt.Sprintf("service name. Mandatory. %s", waiter.NamesMessage()),
				Required: true,
			},
			&cli.IntFlag{
				Name:    "timeout",
				Aliases: []string{"t"},
				Usage:   "timeout in seconds between repeats",
				Value:   int(waiter.DefaultTimeout.Seconds()),
			},
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"l"},
				Usage:   "number of repeats",
				Value:   waiter.DefaultLimit,
			},
		},
		Action: func(c *cli.Context) error {
			service := c.String(serviceFlagName)
			timeout := c.Int(timeoutFlagName)
			limit := c.Int(limitFlagName)

			start := time.Now()

			ready, err := waiter.Wait(service, timeout, limit)
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
