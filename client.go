package main

import (
	"bytes"
	"net"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func doConn(proto, addr string) {
	logBuf := new(bytes.Buffer)
	bufLogger := zerolog.New(logBuf).With().Timestamp().Logger()

	t0 := time.Now()
	conn, err := net.DialTimeout(proto, addr, time.Millisecond*800)
	if err != nil {
		log.Info().
			Str("proto", proto).
			Str("event", "dial").
			Err(err).
			Msg("dial error")
		return
	}

	// Close conn on return.
	defer conn.Close()

	// Set read and write timeout.
	conn.SetDeadline(time.Now().Add(time.Millisecond * 800))

	bufLogger.Info().
		Str("local_addr", conn.LocalAddr().String()).
		Str("remote_addr", conn.RemoteAddr().String()).
		TimeDiff("uptime_seconds", time.Now(), tInit).
		Str("ecs_task_arn", ecsTaskARN).
		Msg("echo request")

	writeBuf := logBuf.Bytes()

	t1 := time.Now()
	writeBytesCount, err := conn.Write(writeBuf)
	if err != nil {
		log.Info().
			Str("proto", proto).
			Str("event", "write").
			Err(err).
			Msg("write error")
		return
	}

	readBuf := make([]byte, 4096)

	t2 := time.Now()
	readBytesCount, err := conn.Read(readBuf)

	t3 := time.Now()
	log.Info().
		Str("proto", proto).
		Str("event", "read").
		Bytes("send_payload", writeBuf).
		Int("send_payload_size", writeBytesCount).
		Bytes("recv_payload", readBuf[:readBytesCount]).
		Int("recv_payload_size", readBytesCount).
		Str("local_addr", conn.LocalAddr().String()).
		Str("remote_addr", conn.RemoteAddr().String()).
		TimeDiff("conn_dur_conn", t1, t0).
		TimeDiff("conn_dur_write", t2, t1).
		TimeDiff("conn_dur_read", t3, t2).
		TimeDiff("conn_dur_total", t3, t1).
		Str("ecs_task_arn", ecsTaskARN).
		TimeDiff("uptime_seconds", time.Now(), tInit).
		Msg("connection read")
}

func client(proto, addr string, interval time.Duration) {
	ticker := time.NewTicker(1000 * time.Millisecond)
	for ; true; <-ticker.C {
		go doConn(proto, addr)
	}
}
