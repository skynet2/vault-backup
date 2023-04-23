package main

import (
	"os"

	"github.com/prometheus/client_golang/prometheus/push"
)

func pushMetrics(prometheusPushGatewayUrl string) error {
	if prometheusPushGatewayUrl == "" {
		return nil
	}

	prometheusJobName := "default_job"
	if v := os.Getenv("PROMETHEUS_JOB_NAME"); v != "" {
		prometheusJobName = v
	}

	return push.New(prometheusPushGatewayUrl, prometheusJobName).
		Collector(successCounter).
		Collector(failCounter).
		Push()
}
