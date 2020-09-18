package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	defaultAppMode              = "client"
	defaultPort                 = "2222"
	defaultClientProto          = "tcp"
	defaultClientTargetHost     = "127.0.0.1"
	defaultClientConnIntervalMS = "1000"
	defaultClientConnCount      = "0"
)

type ecsTaskMetaData struct {
	TaskARN string
}

var tInit time.Time
var outboundIP net.IP
var ecsTaskARN string

func main() {
	// Get the current time, for measuring uptime.
	tInit = time.Now()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.DurationFieldUnit = time.Second
	log.Info().
		Str("event", "start").
		TimeDiff("uptime_seconds", time.Now(), tInit).
		Msg("starting")

	// Get ECS metadata.
	getECSMetadata()

	// Get outbound IP address.
	var err error
	outboundIP, err = getOutboundIP()
	if err != nil {
		log.Info().
			Str("event", "get_outbound_ip").
			Err(err).
			Msg("unable to get outbound IP")
	}

	appMode, ok := os.LookupEnv("APP_MODE")
	if !ok {
		appMode = defaultAppMode
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}

	clientProto, ok := os.LookupEnv("CLIENT_PROTO")
	if !ok {
		clientProto = defaultClientProto
	}

	clientTargetHost, ok := os.LookupEnv("CLIENT_TARGET_HOST")
	if !ok {
		clientTargetHost = defaultClientTargetHost
	}

	clientConnIntervalMSStr, ok := os.LookupEnv("CLIENT_CONN_INTERVAL_MS")
	if !ok {
		clientConnIntervalMSStr = defaultClientConnIntervalMS
	}
	clientConnIntervalMS, _ := strconv.Atoi(clientConnIntervalMSStr)
	clientConnInterval := time.Millisecond * time.Duration(clientConnIntervalMS)

	switch appMode {
	case "client":
		go client(
			clientProto,
			fmt.Sprintf("%s:%s", clientTargetHost, port),
			clientConnInterval,
		)

	case "server":
		// Start UDP and TCP echo servers.
		go serveEchoUDP(fmt.Sprintf(":%s", port))
		go serveEchoTCP(fmt.Sprintf(":%s", port))

	default:
		log.Error().
			Str("app_mode", appMode).
			Msg("unknown app mode")
		os.Exit(1)
	}

	_, ok = os.LookupEnv("HANDLE_SIGNALS")
	if ok {
		go handleSigs()
	}

	select {}
}
