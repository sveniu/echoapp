package main

import (
	"bytes"
	"net"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func serveEchoUDP(listenAddr string) {
	logBuf := new(bytes.Buffer)
	bufLogger := zerolog.New(logBuf).With().Timestamp().Logger()

	pc, err := net.ListenPacket("udp", listenAddr)
	if err != nil {
		log.Info().
			Str("proto", "UDP").
			Str("event", "listenpacket").
			Err(err).
			Msg("listenpacket error")
		return
	}

	// Close packet conn on return.
	defer pc.Close()

	log.Info().
		Str("proto", "UDP").
		Str("event", "listenpacket").
		Str("listen_address", listenAddr).
		Msg("listener started")

	for {
		buf := make([]byte, 4096)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			log.Info().
				Str("proto", "UDP").
				Str("event", "readfrom").
				Str("local_addr", pc.LocalAddr().String()).
				Str("remote_addr", addr.String()).
				Err(err).
				Msg("read error")
			continue
		}
		log.Info().
			Str("proto", "UDP").
			Str("event", "readfrom").
			Str("local_addr", pc.LocalAddr().String()).
			Str("remote_addr", addr.String()).
			Int("recv_size", n).
			Bytes("recv_payload", buf[:n]).
			TimeDiff("uptime_seconds", time.Now(), tInit).
			Str("ecs_task_arn", ecsTaskARN).
			Msg("read data")
		go func(pc net.PacketConn, addr net.Addr, buf []byte) {
			bufLogger.Info().
				IPAddr("local_outbound_ip", outboundIP).
				Str("remote_addr", addr.String()).
				Bytes("echo_request", buf).
				TimeDiff("uptime_seconds", time.Now(), tInit).
				Str("ecs_task_arn", ecsTaskARN).
				Msg("echo reply")
			pc.WriteTo(logBuf.Bytes(), addr)
			logBuf.Reset()
		}(pc, addr, buf[:n])
	}
}

func serveEchoTCP(listenAddr string) {
	logBuf := new(bytes.Buffer)
	bufLogger := zerolog.New(logBuf).With().Timestamp().Logger()

	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Info().
			Str("proto", "TCP").
			Str("event", "listen").
			Err(err).
			Msg("listen error")
		return
	}

	// Close listener on return.
	defer l.Close()

	log.Info().
		Str("proto", "TCP").
		Str("event", "listen").
		Str("listen_address", listenAddr).
		Msg("listener started")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Info().
				Str("proto", "TCP").
				Str("event", "accept").
				Err(err).
				Msg("accept error")
			continue
		}
		t0 := time.Now()

		// Close conn on return.
		defer conn.Close()

		// Set read and write timeout.
		conn.SetDeadline(time.Now().Add(time.Millisecond * 800))

		log.Info().
			Str("proto", "TCP").
			Str("event", "accept").
			Str("local", conn.LocalAddr().String()).
			Str("remote_addr", conn.RemoteAddr().String()).
			TimeDiff("uptime_seconds", time.Now(), tInit).
			Msg("got connection")

		go func(conn net.Conn, startTime time.Time) {
			for {
				buf := make([]byte, 4096)
				n, err := conn.Read(buf)
				if err != nil {
					if err.Error() == "EOF" {
						log.Info().
							Str("proto", "TCP").
							Str("event", "close").
							Str("local_addr", conn.LocalAddr().String()).
							Str("remote_addr", conn.RemoteAddr().String()).
							TimeDiff("uptime_seconds", time.Now(), tInit).
							Str("ecs_task_arn", ecsTaskARN).
							TimeDiff("dur_accept_close", time.Now(), t0).
							Msg("connection closed")
						return
					}
					log.Info().
						Str("proto", "TCP").
						Str("event", "read").
						Str("local_addr", conn.LocalAddr().String()).
						Str("remote_addr", conn.RemoteAddr().String()).
						Err(err).
						Msg("read error")
					return
				}
				log.Info().
					Str("proto", "TCP").
					Str("event", "read").
					Str("local_addr", conn.LocalAddr().String()).
					Str("remote_addr", conn.RemoteAddr().String()).
					Int("recv_size", n).
					Bytes("recv_payload", buf[:n]).
					TimeDiff("uptime_seconds", time.Now(), tInit).
					Str("ecs_task_arn", ecsTaskARN).
					TimeDiff("dur_accept_read", time.Now(), t0).
					Msg("read data")

				bufLogger.Info().
					IPAddr("local_outbound_ip", outboundIP).
					Str("remote_addr", conn.RemoteAddr().String()).
					Bytes("echo_request", buf[:n]).
					TimeDiff("uptime_seconds", time.Now(), tInit).
					Str("ecs_task_arn", ecsTaskARN).
					Msg("echo reply")
				conn.Write(logBuf.Bytes())
				logBuf.Reset()
			}
		}(conn, t0)
	}
}
