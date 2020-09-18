package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/sys/unix"
)

func handleSigs() {
	// Handle all signals.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	for {
		sig := <-sigs
		sigName := unix.SignalName(sig.(syscall.Signal))
		sigNum := unix.SignalNum(sigName)
		log.Info().
			Str("event", "signal").
			Str("sig", sigName).
			Int("signum", int(sigNum)).
			TimeDiff("uptime_seconds", time.Now(), tInit).
			Msg("caught signal")

		if sig == unix.SIGTERM || sig == unix.SIGINT {
			go func() {
				tIntr := time.Now()
				ticker := time.NewTicker(200 * time.Millisecond)
				for ; true; <-ticker.C {
					log.Info().
						Str("event", "limbo").
						TimeDiff("uptime_seconds", time.Now(), tInit).
						TimeDiff("afloat_seconds", time.Now(), tIntr).
						Str("ecs_task_arn", ecsTaskARN).
						Msg("staying alive")
				}
			}()
		}
	}
}
