package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "1.1.1.1:53")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

func getECSMetadata() {
	key := "ECS_CONTAINER_METADATA_URI_V3"
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Info().
			Str("var_name", key).
			Msg("env var not set")
		return
	}

	url := fmt.Sprintf("%s/task", val)
	r, err := http.Get(url)
	if err != nil {
		log.Info().
			Str("url", url).
			Err(err).
			Msg("error while fetching http url")
		return
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		log.Info().
			Int("http_status_code_actual", r.StatusCode).
			Int("http_status_code_expected", http.StatusOK).
			Msg("unexpected http response status code")
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Info().
			Err(err).
			Msg("failed to read http response body")
		return
	}
	log.Info().
		// Remove carriage returns, ref
		// https://godoc.org/github.com/rs/zerolog#Context.RawJSON
		RawJSON("ecs_metadata", bytes.ReplaceAll(bodyBytes, []byte("\r"), []byte{})).
		Msg("got ecs metadata")

	ecsMeta := &ecsTaskMetaData{}
	if err := json.Unmarshal(bodyBytes, ecsMeta); err != nil {
		log.Info().
			Err(err).
			Msg("failed to unmarshal json")
		return
	}

	ecsTaskARN = ecsMeta.TaskARN
	log.Info().
		Str("task_arn", ecsTaskARN).
		Msg("found ecs task arn")
}
